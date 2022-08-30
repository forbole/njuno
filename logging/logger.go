package logging

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	tmctypes "github.com/tendermint/tendermint/rpc/core/types"

	"github.com/forbole/njuno/modules"
	"github.com/forbole/njuno/types"
)

const (
	LogKeyModule  = "module"
	LogKeyHeight  = "height"
	LogKeyTxHash  = "tx_hash"
	LogKeyMsgType = "msg_type"
)

// Logger defines a function that takes an error and logs it.
type Logger interface {
	SetLogLevel(level string) error
	SetLogFormat(format string) error

	Info(msg string, keyvals ...interface{})
	Debug(msg string, keyvals ...interface{})
	Error(msg string, keyvals ...interface{})

	GenesisError(module modules.Module, err error)
	BlockError(module modules.Module, block *tmctypes.ResultBlock, err error)
	EventsError(module modules.Module, results *tmctypes.ResultBlock, err error)
	TxError(module modules.Module, tx *types.TxResponse, err error)
	MsgError(module modules.Module, tx *types.TxResponse, msg sdk.Msg, err error)
}
