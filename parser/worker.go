package parser

import (
	"encoding/json"
	"fmt"

	"github.com/MonikaCat/njuno/logging"

	"github.com/cosmos/cosmos-sdk/codec"

	"github.com/MonikaCat/njuno/database"
	"github.com/MonikaCat/njuno/types/config"

	"github.com/MonikaCat/njuno/modules"

	sdk "github.com/cosmos/cosmos-sdk/types"
	tmctypes "github.com/tendermint/tendermint/rpc/core/types"
	tmtypes "github.com/tendermint/tendermint/types"

	"github.com/MonikaCat/njuno/node"
	"github.com/MonikaCat/njuno/types"
	bdtypes "github.com/MonikaCat/njuno/types"
	"github.com/MonikaCat/njuno/types/utils"
)

// Worker defines a job consumer that is responsible for getting and
// aggregating block and associated data and exporting it to a database.
type Worker struct {
	index int

	queue   types.HeightQueue
	codec   codec.BinaryMarshaler
	modules []modules.Module

	node           node.Node
	db             database.Database
	logger         logging.Logger
	validatorsList *types.ValidatorsList
}

// NewWorker allows to create a new Worker implementation.
func NewWorker(ctx *Context, queue types.HeightQueue, index int) Worker {
	return Worker{
		index:          index,
		codec:          ctx.EncodingConfig.Marshaler,
		node:           ctx.Node,
		queue:          queue,
		db:             ctx.Database,
		modules:        ctx.Modules,
		logger:         ctx.Logger,
		validatorsList: ctx.ValidatorsList,
	}
}

// Start starts a worker by listening for new jobs (block heights) from the
// given worker queue. Any failed job is logged and re-enqueued.
func (w Worker) Start() {
	logging.WorkerCount.Inc()
	nodeInfo, err := w.node.Genesis()
	if err != nil {
		w.logger.Error("error while getting genesis info from the node ", "err", err)
	}

	for i := range w.queue {
		if err := w.ProcessIfNotExists(i); err != nil {
			// re-enqueue any failed job
			go func() {
				w.logger.Error("re-enqueueing failed block", "height", i, "err", err)
				w.queue <- i
			}()
		}

		logging.WorkerHeight.WithLabelValues(fmt.Sprintf("%d", w.index), nodeInfo.Genesis.ChainID).Set(float64(i))
	}
}

// ProcessIfNotExists defines the job consumer workflow. It will fetch a block for a given
// height and associated metadata and export it to a database if it does not exist yet.
// It returns an error if any export process fails.
func (w Worker) ProcessIfNotExists(height int64) error {
	exists, err := w.db.HasBlock(height)
	if err != nil {
		return fmt.Errorf("error while searching for block: %s", err)
	}

	if exists {
		w.logger.Debug("skipping already exported block", "height", height)
		return nil
	}

	return w.Process(height)
}

// Process fetches a block for a given height and associated metadata and export it to a database.
// It returns an error if any export process fails.
func (w Worker) Process(height int64) error {
	if height == 0 {
		cfg := config.Cfg.Parser

		genesisDoc, genesisState, err := utils.GetGenesisDocAndState(cfg.GenesisFilePath, w.node)
		if err != nil {
			return fmt.Errorf("failed to get genesis: %s", err)
		}

		return w.HandleGenesis(genesisDoc, genesisState)
	}

	w.logger.Debug("processing block", "height", height)

	block, err := w.node.Block(height)
	if err != nil {
		return fmt.Errorf("failed to get block from node: %s", err)
	}

	events, err := w.node.BlockResults(height)
	if err != nil {
		return fmt.Errorf("failed to get block results from node: %s", err)
	}

	txs, err := w.UnmarshalTxs(block)
	if err != nil {
		return fmt.Errorf("failed to get transactions for block: %s", err)
	}

	vals, err := w.node.Validators(height)
	if err != nil {
		return fmt.Errorf("failed to get validators for block: %s", err)
	}

	return w.ExportBlock(block, events, txs, vals)
}

// ProcessTransactions fetches transactions for a given height and stores them into the database.
// It returns an error if the export process fails.
func (w Worker) ProcessTransactions(height int64) error {
	block, err := w.node.Block(height)
	if err != nil {
		return fmt.Errorf("failed to get block from node: %s", err)
	}

	txs, err := w.UnmarshalTxs(block)
	if err != nil {
		return fmt.Errorf("failed to get transactions for block: %s", err)
	}

	return w.ExportTxs(txs)
}

// HandleGenesis accepts a GenesisDoc and calls all the registered genesis handlers
// in the order in which they have been registered.
func (w Worker) HandleGenesis(genesisDoc *tmtypes.GenesisDoc, appState map[string]json.RawMessage) error {
	// Call the genesis handlers
	for _, module := range w.modules {
		if genesisModule, ok := module.(modules.GenesisModule); ok {
			if err := genesisModule.HandleGenesis(genesisDoc, appState); err != nil {
				w.logger.GenesisError(module, err)
			}
		}
	}

	return nil
}

// SaveValidators persists a list of Tendermint validators with an address and a
// consensus public key. An error is returned if the public key cannot be Bech32
// encoded or if the DB write fails.
func (w Worker) SaveValidators(vals []*tmtypes.Validator, height int64) error {
	var validators []types.Validator
	var validatorsDesc []types.ValidatorDescription

	for _, val := range vals {
		consAddr := sdk.ConsAddress(val.Address).String()

		consPubKey, err := types.ConvertValidatorPubKeyToBech32String(val.PubKey)
		if err != nil {
			return fmt.Errorf("failed to convert validator public key for validators %s: %s", consAddr, err)
		}

		validatorAddress, err := sdk.ValAddressFromHex(val.Address.String())
		if err != nil {
			fmt.Printf("failed to convert validator address from hex: %s", err)
		}

		for _, entry := range w.validatorsList.Validators {
			// compare with address from yaml file
			if entry.Validator.Hex == val.Address.String() {
				// store validators
				validators = append(validators, types.NewValidator(consAddr, validatorAddress.String(), consPubKey, entry.Validator.Address, height))

				// store validators description
				validatorsDesc = append(validatorsDesc, types.NewValidatorDescription(consAddr, entry.Validator.Details, entry.Validator.Identity, entry.Validator.Moniker, height))
			}
		}

	}

	err := w.db.SaveValidators(validators)
	if err != nil {
		return fmt.Errorf("error while saving validators: %s", err)
	}

	err = w.db.SaveValidatorDescription(validatorsDesc)
	if err != nil {
		return fmt.Errorf("error while saving validators description: %s", err)
	}

	return nil
}

// ExportBlock accepts a finalized block and a corresponding set of transactions
// and persists them to the database along with attributable metadata. An error
// is returned if the write fails.
func (w Worker) ExportBlock(
	b *tmctypes.ResultBlock, r *tmctypes.ResultBlockResults, txs []types.TxResponse, vals *tmctypes.ResultValidators,
) error {

	// Save all validators
	err := w.SaveValidators(vals.Validators, b.Block.Height)
	if err != nil {
		return err
	}

	// Make sure the proposer exists
	proposerAddr := sdk.ConsAddress(b.Block.ProposerAddress)
	val := findValidatorByAddr(proposerAddr.String(), vals)
	if val == nil {
		return fmt.Errorf("failed to find validator by proposer address %s: %s", proposerAddr.String(), err)
	}

	// Save the block
	err = w.db.SaveBlock(types.NewBlockFromTmBlock(b, 0))
	if err != nil {
		return fmt.Errorf("failed to persist block: %s", err)
	}

	// Save the commits
	err = w.ExportCommit(b.Block.LastCommit, vals)
	if err != nil {
		return err
	}

	// Call the block handlers
	for _, module := range w.modules {
		if blockModule, ok := module.(modules.BlockModule); ok {
			err = blockModule.HandleBlock(b, r, vals)
			if err != nil {
				w.logger.BlockError(module, b, err)
			}
		}
	}

	// Export the transactions
	return w.ExportTxs(txs)
}

// ExportCommit accepts a block commitment and a corresponding set of
// validators for the commitment and persists them to the database. An error is
// returned if any write fails or if there is any missing aggregated data.
func (w Worker) ExportCommit(commit *tmtypes.Commit, vals *tmctypes.ResultValidators) error {
	var signatures []*types.CommitSig

	for _, commitSig := range commit.Signatures {
		// Avoid empty commits
		if commitSig.Signature == nil {
			continue
		}

		valAddr := sdk.ConsAddress(commitSig.ValidatorAddress)
		val := findValidatorByAddr(valAddr.String(), vals)
		if val == nil {
			return fmt.Errorf("failed to find validator by commit validator address %s", valAddr.String())
		}

		signatures = append(signatures, types.NewCommitSig(
			types.ConvertValidatorAddressToBech32String(commitSig.ValidatorAddress),
			val.VotingPower,
			val.ProposerPriority,
			commit.Height,
			commitSig.Timestamp,
		))
	}

	err := w.db.SaveCommitSignatures(signatures)
	if err != nil {
		return fmt.Errorf("error while saving commit signatures: %s", err)
	}

	return nil
}

// ExportTxs accepts a slice of transactions and persists then inside the database.
// An error is returned if the write fails.
func (w Worker) ExportTxs(txs []types.TxResponse) error {
	// Handle all transactions inside the block
	for _, tx := range txs {
		// Save  transaction in database
		err := w.db.SaveTx(tx)
		if err != nil {
			return fmt.Errorf("failed to handle transaction with hash %s: %s", tx.Hash, err)
		}
	}

	return nil
}

// UnmarshalTxs process all transactions contained in a block
func (w Worker) UnmarshalTxs(block *tmctypes.ResultBlock) ([]bdtypes.TxResponse, error) {
	txResponses := make([]bdtypes.TxResponse, len(block.Block.Txs))

	// get tx details from the block
	var transaction bdtypes.TxResponse
	for _, t := range block.Block.Txs {
		err := json.Unmarshal(t, &transaction)
		if err != nil {
			// continue
		}

		txResponses = append(txResponses, bdtypes.NewTxResponse(transaction.Fee, transaction.Memo, transaction.Msg, transaction.Signatures, fmt.Sprintf("%X", t.Hash()), block.Block.Height))
	}

	return txResponses, nil
}
