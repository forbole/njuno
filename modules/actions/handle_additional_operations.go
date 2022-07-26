package actions

import (
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/forbole/njuno/modules/actions/handlers"
	actionstypes "github.com/forbole/njuno/modules/actions/types"
)

var (
	waitGroup sync.WaitGroup
)

func (m *Module) RunAdditionalOperations() error {
	// Build the worker
	context := actionstypes.NewContext(m.node)
	worker := actionstypes.NewActionsWorker(context)

	// -- Register the Account Balance endpoint --
	worker.RegisterHandler("/account_balance", handlers.AccountBalanceHandler)

	// -- Staking Delegator --
	worker.RegisterHandler("/delegation_total", handlers.TotalDelegationsAmountHandler)

	// Listen for and trap any OS signal to gracefully shutdown and exit
	m.trapSignal()

	// Start the worker
	waitGroup.Add(1)
	go worker.Start(m.cfg.Port)

	// Block main process (signal capture will call WaitGroup's Done)
	waitGroup.Wait()
	return nil
}

// trapSignal will listen for any OS signal and invoke Done on the main
// WaitGroup allowing the main process to gracefully exit.
func (m *Module) trapSignal() {
	var sigCh = make(chan os.Signal, 1)

	signal.Notify(sigCh, syscall.SIGTERM)
	signal.Notify(sigCh, syscall.SIGINT)

	go func() {
		defer m.node.Stop()
		defer waitGroup.Done()
	}()
}
