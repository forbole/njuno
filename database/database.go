package database

import (
	"time"

	"github.com/cosmos/cosmos-sdk/simapp/params"

	"github.com/forbole/njuno/logging"

	databaseconfig "github.com/forbole/njuno/database/config"

	sdk "github.com/cosmos/cosmos-sdk/types"
	dbtypes "github.com/forbole/njuno/database/types"
	"github.com/forbole/njuno/types"
)

// Database represents an abstract database that can be used to save data inside it
type Database interface {
	// Close closes the connection to the database
	Close()

	// GetBlockHeightTimeDayAgo returns block height from day ago.
	// An error is returned if the operation fails.
	GetBlockHeightTimeDayAgo(now time.Time) (dbtypes.BlockRow, error)

	// GetBlockHeightTimeHourAgo returns block height from one hour ago.
	// An error is returned if the operation fails.
	GetBlockHeightTimeHourAgo(now time.Time) (dbtypes.BlockRow, error)

	// GetBlockHeightTimeMinuteAgo returns block height from one minute ago.
	// An error is returned if the operation fails.
	GetBlockHeightTimeMinuteAgo(now time.Time) (dbtypes.BlockRow, error)

	// GetGenesis returns the genesis details.
	// An error is returned if the operation fails.
	GetGenesis() (*types.Genesis, error)

	// GetLastBlock returns the latest block store in database.
	// An error is returned if the operation fails.
	GetLastBlock() (*dbtypes.BlockRow, error)

	// GetLastBlockHeight returns the latest block height stored in database.
	// An error is returned if the operation fails.
	GetLastBlockHeight() (int64, error)

	// GetValidatorsDescription returns validators description stored in database.
	// An error is returned if the operation fails.
	GetValidatorsDescription() ([]types.ValidatorDescription, error)

	// GetTokensPriceID returns token ID stored in database.
	// An error is returned if the operation fails.
	GetTokensPriceID() ([]string, error)

	// HasBlock tells whether or not the database has already stored the block having the given height.
	// An error is returned if the operation fails.
	HasBlock(height int64) (bool, error)

	// SaveAverageBlockTimeGenesis stores the average
	// block time from genesis.
	// An error is returned if the operation fails.
	SaveAverageBlockTimeGenesis(averageTime float64, height int64) error

	// SaveAverageBlockTimePerDay stores the average
	// block time per day.
	// An error is returned if the operation fails.
	SaveAverageBlockTimePerDay(averageTime float64, height int64) error

	// SaveAverageBlockTimePerHour stores the average
	// block time per hour.
	// An error is returned if the operation fails.
	SaveAverageBlockTimePerHour(averageTime float64, height int64) error

	// SaveAverageBlockTimePerMin stores the average
	// block time per minute.
	// An error is returned if the operation fails.
	SaveAverageBlockTimePerMin(averageTime float64, height int64) error

	// SaveBlock will be called when a new block is parsed, passing the block itself
	// and the transactions contained inside that block.
	// An error is returned if the operation fails.
	// NOTE. For each transaction inside txs, SaveTx will be called as well.
	SaveBlock(block *types.Block) error

	// SaveCommitSignatures stores a  slice of validator commit signatures.
	// An error is returned if the operation fails.
	SaveCommitSignatures(signatures []*types.CommitSig) error

	// SaveDoubleSignEvidence stores double sign record in database.
	// An error is returned if the operation fails.
	SaveDoubleSignEvidence(evidence types.DoubleSignEvidence) error

	// SaveGenesis stores the genesis details in database.
	// An error is returned if the operation fails.
	SaveGenesis(genesis *types.Genesis) error

	// SaveIBCTransferParams stores the ibc transfer params value in database.
	// An error is returned if the operation fails.
	SaveIBCTransferParams(params *types.IBCTransferParams) error

	// SaveInflation stores the inflation value in database.
	// An error is returned if the operation fails.
	SaveInflation(inflation string, height int64) error

	// SaveStakingPool stores the staking pool value in database.
	// An error is returned if the operation fails.
	SaveStakingPool(pool *types.StakingPool) error

	// SaveSupply stores a total supply value in database.
	// An error is returned if the operation fails.
	SaveSupply(coins sdk.Coins, height int64) error

	// SaveToken stores the token details in database.
	// An error is returned if the operation fails.
	SaveToken(token types.Token) error

	// SaveTokensPrice stores tokens price in database.
	// An error is returned if the operation fails.
	SaveTokensPrice(prices []types.TokenPrice) error

	// SaveTx stores transaction contained inside a block in database.
	// An error is returned if the operation fails.
	SaveTx(tx types.TxResponse) error

	// SaveValidators stores a list of validators in database.
	// An error is returned if the operation fails.
	SaveValidators(validators []types.Validator) error

	// SaveValidatorCommission stores validators commission value in database.
	// An error is returned if the operation fails.
	SaveValidatorCommission(data []types.ValidatorCommission) error

	// SaveValidatorDescription stores the validators description in database.
	// An error is returned if the operation fails.
	SaveValidatorDescription(description []types.ValidatorDescription) error

	// SaveValidatorsStatus stores the validators status in database.
	// An error is returned if the operation fails.
	SaveValidatorsStatus(validatorsStatus []types.ValidatorStatus) error

	// SaveValidatorsVotingPower stores a list of validators voting power in database.
	// An error is returned if the operation fails.
	SaveValidatorsVotingPower(entries []types.ValidatorVotingPower) error
}

// PruningDb represents a database that supports pruning properly
type PruningDb interface {
	// GetLastPruned returns the last height at which the database was pruned
	GetLastPruned() (int64, error)

	// Prune prunes the data for the given height, returning any error
	Prune(height int64) error

	// StoreLastPruned saves the last height at which the database was pruned
	StoreLastPruned(height int64) error
}

// Context contains the data that might be used to build a Database instance
type Context struct {
	Cfg            databaseconfig.Config
	EncodingConfig *params.EncodingConfig
	Logger         logging.Logger
}

// NewContext allows to build a new Context instance
func NewContext(cfg databaseconfig.Config, encodingConfig *params.EncodingConfig, logger logging.Logger) *Context {
	return &Context{
		Cfg:            cfg,
		EncodingConfig: encodingConfig,
		Logger:         logger,
	}
}

// Builder represents a method that allows to build any database from a given codec and configuration
type Builder func(ctx *Context) (Database, error)
