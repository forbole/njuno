package database

import (
	"time"

	"github.com/cosmos/cosmos-sdk/simapp/params"

	"github.com/MonikaCat/njuno/logging"

	databaseconfig "github.com/MonikaCat/njuno/database/config"

	dbtypes "github.com/MonikaCat/njuno/database/types"
	"github.com/MonikaCat/njuno/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Database represents an abstract database that can be used to save data inside it
type Database interface {
	// HasBlock tells whether or not the database has already stored the block having the given height.
	// An error is returned if the operation fails.
	HasBlock(height int64) (bool, error)

	// SaveBlock will be called when a new block is parsed, passing the block itself
	// and the transactions contained inside that block.
	// An error is returned if the operation fails.
	// NOTE. For each transaction inside txs, SaveTx will be called as well.
	SaveBlock(block *types.Block) error

	// SaveTx will be called to save each transaction contained inside a block.
	// An error is returned if the operation fails.
	SaveTx(tx types.TxResponse) error

	// SaveValidators stores a list of validators if they do not already exist.
	// An error is returned if the operation fails.
	SaveValidators(validators []types.Validator) error

	// SaveValidatorsVotingPowers stores a list of validators voting power.
	// An error is returned if the operation fails.
	SaveValidatorsVotingPowers(entries []types.ValidatorVotingPower) error

	// SaveValidatorDescription stores the validators description.
	// An error is returned if the operation fails.
	SaveValidatorDescription(description []types.ValidatorDescription) error

	// SaveCommitSignatures stores a  slice of validator commit signatures.
	// An error is returned if the operation fails.
	SaveCommitSignatures(signatures []*types.CommitSig) error

	// SaveMessage stores a single message.
	// An error is returned if the operation fails.
	SaveMessage(msg *types.Message) error

	// SaveSupply stores a total supply value.
	// An error is returned if the operation fails.
	SaveSupply(coins sdk.Coins, height int64) error

	// SaveInflation stores the inflation value.
	// An error is returned if the operation fails.
	SaveInflation(inflation string, height int64) error

	// SaveStakingPool stores the staking pool value.
	// An error is returned if the operation fails.
	SaveStakingPool(pool *types.StakingPool) error

	// GetLastBlockHeight returns the latest block height stored in database.
	// An error is returned if the operation fails.
	GetLastBlockHeight() (int64, error)

	// GetLastBlock returns the latest block store in db.
	// An error is returned if the operation fails.
	GetLastBlock() (*dbtypes.BlockRow, error)

	// SaveIBCParams stores the ibc tx params value.
	// An error is returned if the operation fails.
	SaveIBCParams(params *types.IBCParams) error

	// SaveAccountBalance stores the account balance value.
	// An error is returned if the operation fails.
	SaveAccountBalance(balances types.AccountBalance) error

	// SaveToken stores the token details.
	// An error is returned if the operation fails.
	SaveToken(token types.Token) error

	// SaveGenesis stores the genesis details.
	// An error is returned if the operation fails.
	SaveGenesis(genesis *types.Genesis) error

	// GetGenesis returns the genesis details.
	// An error is returned if the operation fails.
	GetGenesis() (*types.Genesis, error)

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

	// GetBlockHeightTimeDayAgo returns block height from day ago.
	// An error is returned if the operation fails.
	GetBlockHeightTimeDayAgo(now time.Time) (dbtypes.BlockRow, error)

	// GetBlockHeightTimeHourAgo returns block height from one hour ago.
	// An error is returned if the operation fails.
	GetBlockHeightTimeHourAgo(now time.Time) (dbtypes.BlockRow, error)

	// GetBlockHeightTimeMinuteAgo returns block height from one minute ago.
	// An error is returned if the operation fails.
	GetBlockHeightTimeMinuteAgo(now time.Time) (dbtypes.BlockRow, error)

	// Close closes the connection to the database
	Close()
}

// PruningDb represents a database that supports pruning properly
type PruningDb interface {
	// Prune prunes the data for the given height, returning any error
	Prune(height int64) error

	// StoreLastPruned saves the last height at which the database was pruned
	StoreLastPruned(height int64) error

	// GetLastPruned returns the last height at which the database was pruned
	GetLastPruned() (int64, error)
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