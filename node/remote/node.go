package remote

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	tmtypes "github.com/tendermint/tendermint/types"

	"github.com/cosmos/cosmos-sdk/codec"

	constypes "github.com/tendermint/tendermint/consensus/types"
	tmjson "github.com/tendermint/tendermint/libs/json"

	"github.com/MonikaCat/njuno/node"
	"github.com/MonikaCat/njuno/types"

	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	bdtypes "github.com/MonikaCat/njuno/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	httpclient "github.com/tendermint/tendermint/rpc/client/http"
	tmctypes "github.com/tendermint/tendermint/rpc/core/types"
	jsonrpcclient "github.com/tendermint/tendermint/rpc/jsonrpc/client"
)

var (
	_ node.Node = &Node{}
)

var nomicNode = "https://app.nomic.io:8443/"

// Node implements a wrapper around both a Tendermint RPCConfig client and a
// chain SDK REST client that allows for essential data queries.
type Node struct {
	ctx        context.Context
	codec      codec.Marshaler
	client     *httpclient.HTTP
	clientNode string // Full (REST client) node

	// rpcClient rpcclient.Client // Tendermint (RPC client) node
	// txServiceClient tx.C
	// grpcConnection  *grpc.ClientConn
}

// NewNode allows to build a new Node instance
func NewNode(cfg *Details, codec codec.Marshaler) (*Node, error) {
	clientNode := "http://138.197.71.46:26657"
	httpClient, err := jsonrpcclient.DefaultHTTPClient(cfg.RPC.Address)
	if err != nil {
		return nil, err
	}

	// Tweak the transport
	httpTransport, ok := (httpClient.Transport).(*http.Transport)
	if !ok {
		return nil, fmt.Errorf("invalid HTTP Transport: %T", httpTransport)
	}
	httpTransport.MaxConnsPerHost = cfg.RPC.MaxConnections

	rpcClient, err := httpclient.NewWithClient(cfg.RPC.Address, "/websocket", httpClient)
	if err != nil {
		return nil, err
	}

	err = rpcClient.Start()
	if err != nil {
		return nil, err
	}

	return &Node{
		ctx:   context.Background(),
		codec: codec,

		client:     rpcClient,
		clientNode: clientNode,
	}, nil
}

// Genesis implements node.Node
func (cp *Node) Genesis() (*tmctypes.ResultGenesis, error) {
	res, err := cp.client.Genesis(cp.ctx)
	if err != nil && strings.Contains(err.Error(), "use the genesis_chunked API instead") {
		return cp.getGenesisChunked()
	}
	return res, err
}

// getGenesisChunked gets the genesis data using the chinked API instead
func (cp *Node) getGenesisChunked() (*tmctypes.ResultGenesis, error) {
	bz, err := cp.getGenesisChunksStartingFrom(0)
	if err != nil {
		return nil, err
	}

	var genDoc *tmtypes.GenesisDoc
	err = tmjson.Unmarshal(bz, &genDoc)
	if err != nil {
		return nil, err
	}

	return &tmctypes.ResultGenesis{Genesis: genDoc}, nil
}

// getGenesisChunksStartingFrom returns all the genesis chunks data starting from the chunk with the given id
func (cp *Node) getGenesisChunksStartingFrom(id uint) ([]byte, error) {
	res, err := cp.client.GenesisChunked(cp.ctx, id)
	if err != nil {
		return nil, fmt.Errorf("error while getting genesis chunk %d out of %d", id, res.TotalChunks)
	}

	bz, err := base64.StdEncoding.DecodeString(res.Data)
	if err != nil {
		return nil, fmt.Errorf("error while decoding genesis chunk %d out of %d", id, res.TotalChunks)
	}

	if id == uint(res.TotalChunks-1) {
		return bz, nil
	}

	nextChunk, err := cp.getGenesisChunksStartingFrom(id + 1)
	if err != nil {
		return nil, err
	}

	return append(bz, nextChunk...), nil
}

// ConsensusState implements node.Node
func (cp *Node) ConsensusState() (*constypes.RoundStateSimple, error) {
	state, err := cp.client.ConsensusState(context.Background())
	if err != nil {
		return nil, err
	}

	var data constypes.RoundStateSimple
	err = tmjson.Unmarshal(state.RoundState, &data)
	if err != nil {
		return nil, err
	}
	return &data, nil
}

// LatestHeight implements node.Node
func (cp *Node) LatestHeight() (int64, error) {
	status, err := cp.client.Status(cp.ctx)
	if err != nil {
		return -1, err
	}

	height := status.SyncInfo.LatestBlockHeight
	return height, nil
}

// Validators implements node.Node
func (cp *Node) Validators(height int64) (*tmctypes.ResultValidators, error) {
	vals := &tmctypes.ResultValidators{
		BlockHeight: height,
	}

	page := 1
	perPage := 100 // maximum 100 entries per page
	stop := false
	for !stop {
		result, err := cp.client.Validators(cp.ctx, &height, &page, &perPage)
		if err != nil {
			return nil, err
		}
		vals.Validators = append(vals.Validators, result.Validators...)
		vals.Count += result.Count
		vals.Total = result.Total
		page += 1
		stop = vals.Count == vals.Total
	}

	return vals, nil
}

// Block implements node.Node
func (cp *Node) Block(height int64) (*tmctypes.ResultBlock, error) {
	return cp.client.Block(cp.ctx, &height)
}

// BlockResults implements node.Node
func (cp *Node) BlockResults(height int64) (*tmctypes.ResultBlockResults, error) {
	return cp.client.BlockResults(cp.ctx, &height)
}

// Txs implements node.Node
func (cp *Node) Txs(block *tmctypes.ResultBlock) ([]bdtypes.TxResponse, error) {
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

// // TxSearch implements node.Node
// func (cp *Node) TxSearch(query string, page *int, perPage *int, orderBy string) (*tmctypes.ResultTxSearch, error) {
// 	return cp.client.TxSearch(cp.ctx, query, false, page, perPage, orderBy)
// }

// SubscribeEvents implements node.Node
func (cp *Node) SubscribeEvents(subscriber, query string) (<-chan tmctypes.ResultEvent, context.CancelFunc, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	eventCh, err := cp.client.Subscribe(ctx, subscriber, query)
	return eventCh, cancel, err
}

// SubscribeNewBlocks implements node.Node
func (cp *Node) SubscribeNewBlocks(subscriber string) (<-chan tmctypes.ResultEvent, context.CancelFunc, error) {
	return cp.SubscribeEvents(subscriber, "tm.event = 'NewBlock'")
}

// Stop implements node.Node
func (cp *Node) Stop() {
	err := cp.client.Stop()
	if err != nil {
		panic(fmt.Errorf("error while stopping proxy: %s", err))
	}
}

// Supply implements node.Node
func (cp *Node) Supply() (sdk.Coins, error) {
	resp, err := http.Get(fmt.Sprintf("%s/cosmos/bank/v1beta1/supply/unom", nomicNode))
	if err != nil {
		return sdk.Coins{}, fmt.Errorf("error while getting total supply: %s", err)
	}

	defer resp.Body.Close()

	bz, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return sdk.Coins{}, fmt.Errorf("error while processing total supply: %s", err)
	}

	var supply *banktypes.QuerySupplyOfResponse
	err = json.Unmarshal(bz, &supply)
	if err != nil {
		return sdk.Coins{}, fmt.Errorf("error while unmarshaling supply: %s", err)
	}
	var totalSupply []sdk.Coin
	totalSupply = append(totalSupply, supply.Amount)

	return totalSupply, nil
}

// Inflation implements node.Node
func (cp *Node) Inflation() (string, error) {
	resp, err := http.Get(fmt.Sprintf("%s/cosmos/mint/v1beta1/inflation", nomicNode))
	if err != nil {
		return "", fmt.Errorf("error while getting inflation: %s", err)
	}

	defer resp.Body.Close()

	bz, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error while processing inflation: %s", err)
	}
	var inflation types.InflationResponse
	err = json.Unmarshal(bz, &inflation)
	if err != nil {
		return "", fmt.Errorf("error while unmarshaling inflation: %s", err)
	}

	return string(inflation.Inflation), nil
}

// StkingPool implements node.Node
func (cp *Node) StakingPool() (stakingtypes.Pool, error) {
	resp, err := http.Get(fmt.Sprintf("%s/cosmos/staking/v1beta1/pool", nomicNode))
	if err != nil {
		return stakingtypes.Pool{}, fmt.Errorf("error while getting staking pool: %s", err)
	}

	defer resp.Body.Close()

	bz, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return stakingtypes.Pool{}, fmt.Errorf("error while processing staking pool: %s", err)
	}

	var stakingPool stakingtypes.Pool
	err = json.Unmarshal(bz, &stakingPool)
	if err != nil {
		return stakingtypes.Pool{}, fmt.Errorf("error while unmarshaling staking pool: %s", err)
	}

	return stakingPool, nil
}

// IBCParams implements node.Node
func (cp *Node) IBCParams() (types.IBCTransactionParams, error) {
	resp, err := http.Get(fmt.Sprintf("%s/ibc/apps/transfer/v1/params", nomicNode))
	if err != nil {
		return types.IBCTransactionParams{}, fmt.Errorf("error while getting ibc params: %s", err)
	}

	defer resp.Body.Close()

	bz, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return types.IBCTransactionParams{}, fmt.Errorf("error while processing ibc params: %s", err)
	}

	var stakingPool types.IBCTransferParams
	err = json.Unmarshal(bz, &stakingPool)
	if err != nil {
		return types.IBCTransactionParams{}, fmt.Errorf("error while unmarshaling ibc params: %s", err)
	}

	return stakingPool.Params, nil
}

// AccountBalance implements node.Node
func (cp *Node) AccountBalance(address string) (sdk.Coins, error) {
	resp, err := http.Get(fmt.Sprintf("%s/cosmos/bank/v1beta1/balances/%s", nomicNode, address))
	if err != nil {
		return sdk.Coins{}, fmt.Errorf("error while getting account balance of address %s: %s", address, err)
	}

	defer resp.Body.Close()

	bz, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return sdk.Coins{}, fmt.Errorf("error while processing account balance of address %s: %s", address, err)
	}

	var balance types.QueryAllBalancesResponse
	err = json.Unmarshal(bz, &balance)
	if err != nil {
		return sdk.Coins{}, fmt.Errorf("error while unmarshaling account balance of address %s: %s", address, err)
	}

	return balance.Balances, nil
}