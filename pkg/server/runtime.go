package server

import (
	"context"
	"fmt"
	"net"

	"sync"

	"github.com/pkg/errors"
	"github.com/tevino/abool"
	"github.com/varas/numserver/pkg/errhandler"
	"github.com/varas/numserver/pkg/line"
	"github.com/varas/numserver/pkg/report"
	"github.com/varas/numserver/pkg/repository"
	"github.com/varas/numserver/pkg/result"
)

// runtime context
type runtime struct {
	isUp           *abool.AtomicBool
	stopped        chan struct{}
	errHandle      errhandler.ErrHandler
	cancelListener context.CancelFunc
	cancelHandlers context.CancelFunc
	cancelRunners  context.CancelFunc
	wgHandlers     sync.WaitGroup
	wgDaemons      sync.WaitGroup
}

// passing config on start enables hot config-reloading
func (r *runtime) start(ctx context.Context, c config, errHandle errhandler.ErrHandler) (err error) {
	r.stopped = make(chan struct{})
	r.errHandle = errHandle

	conns := make(chan net.Conn)
	listener, err := NewListener(c.port, conns)
	if err != nil {
		return errors.Wrap(err, "cannot create connection listener")
	}

	// stop runtime in order
	var ctxListener, ctxHandlers, ctxRunners context.Context
	ctxListener, r.cancelListener = context.WithCancel(ctx)
	ctxHandlers, r.cancelHandlers = context.WithCancel(ctx)
	ctxRunners, r.cancelRunners = context.WithCancel(ctx)

	currentReport := &report.Report{}
	numberRepository := repository.NewInMemoryRepository()

	reportRunner := report.NewRunner(c.reportFlushInterval, currentReport)
	resultRunner, err := result.NewRunner(c.logFlushInterval, c.logPath, c.logFlushBatchSize, numberRepository)
	if err != nil {
		return errors.Wrap(err, "cannot create result runner")
	}

	// stop bg jobs: listener and runners
	r.wgDaemons = sync.WaitGroup{}
	r.wgDaemons.Add(3)
	go func() {
		listerErr := listener.Listen(ctxListener)
		// avoid logging connection closed on teardown
		if r.isUp.IsSet() {
			r.errHandle(listerErr)
		}
		r.wgDaemons.Done()
	}()
	go func() {
		r.errHandle(reportRunner.Run(ctxRunners))
		r.wgDaemons.Done()
	}()
	go func() {
		r.errHandle(resultRunner.Run(ctxRunners))
		r.wgDaemons.Done()
	}()

	terminate := make(chan struct{})

	lineValidator, err := line.NewValidator()
	if err != nil {
		return fmt.Errorf("cannot create line validator: %s", err.Error())
	}

	connHandler := newConnHandler(errHandle, lineValidator, numberRepository, currentReport, conns, terminate)

	r.wgHandlers = sync.WaitGroup{}
	r.wgHandlers.Add(c.concurrentClients)
	for w := c.concurrentClients; w > 0; w-- {
		go func() {
			connHandler.run(ctxHandlers)
			r.wgHandlers.Done()
		}()
	}

	go r.waitForClientTermination(terminate)
	go r.waitForContextTermination(ctx)

	r.isUp = abool.NewBool(true)

	return nil
}

func (r *runtime) stop() {
	wasStopped := r.isUp.SetToIf(true, false)
	if !wasStopped {
		return
	}

	r.cancelListener()

	// stop conn handlers
	r.cancelHandlers()
	r.wgHandlers.Wait()

	// stop result & report runners
	r.cancelRunners()
	r.wgDaemons.Wait()

	close(r.stopped)
}

func (r *runtime) waitForClientTermination(termination <-chan struct{}) {
	<-termination
	r.stop()
}

func (r *runtime) waitForContextTermination(ctx context.Context) {
	<-ctx.Done()
	r.stop()
}
