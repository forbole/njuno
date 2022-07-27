package node

import (
	"context"

	"github.com/MonikaCat/njuno/types"
	bdtypes "github.com/MonikaCat/njuno/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	constypes "github.com/tendermint/tendermint/consensus/types"
	tmctypes "github.com/tendermint/tendermint/rpc/core/types"
)

type Node interface {
	// Genesis returns the genesis state
	Genesis() (*tmctypes.ResultGenesis, error)

	// ConsensusState returns the consensus state of the chain
	ConsensusState() (*constypes.RoundStateSimple, error)

	// LatestHeight returns the latest block height on the active chain. An error
	// is returned if the query fails.
	LatestHeight() (int64, error)

	// Validators returns all the known Tendermint validators for a given block
	// height. An error is returned if the query fails.
	Validators(height int64) (*tmctypes.ResultValidators, error)

	// Block queries for a block by height. An error is returned if the query fails.
	Block(height int64) (*tmctypes.ResultBlock, error)

	// BlockResults queries the results of a block by height. An error is returnes if the query fails
	BlockResults(height int64) (*tmctypes.ResultBlockResults, error)

	// Txs queries for all the transactions in a block. Transactions are returned
	// in the sdk.TxResponse format which internally contains an sdk.Tx. An error is
	// returned if any query fails.
	Txs(block *tmctypes.ResultBlock) ([]bdtypes.TxResponse, error)

	// TxSearch defines a method to search for a paginated set of transactions by DeliverTx event search criteria.
	// TxSearch(query string, page *int, perPage *int, orderBy string) (*tmctypes.ResultTxSearch, error)

	// SubscribeEvents subscribes to new events with the given query through the RPCConfig
	// client with the given subscriber name. A receiving only channel, context
	// cancel function and an error is returned. It is up to the caller to cancel
	// the context and handle any errors appropriately.
	SubscribeEvents(subscriber, query string) (<-chan tmctypes.ResultEvent, context.CancelFunc, error)

	// SubscribeNewBlocks subscribes to the new block event handler through the RPCConfig
	// client with the given subscriber name. An receiving only channel, context
	// cancel function and an error is returned. It is up to the caller to cancel
	// the context and handle any errors appropriately.
	SubscribeNewBlocks(subscriber string) (<-chan tmctypes.ResultEvent, context.CancelFunc, error)

	// Stop defers the node stop execution to the client.
	Stop()

	Supply() (sdk.Coins, error)

	Inflation() (string, error)

	StakingPool() (stakingtypes.Pool, error)

	IBCParams() (types.IBCTransactionParams, error)

	AccountBalance(address string) (sdk.Coins, error)
}