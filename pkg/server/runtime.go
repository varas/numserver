package server

import (
	"context"
	"fmt"
	"net"

	"sync"

	"github.com/varas/numserver/pkg/errhandler"
	"github.com/varas/numserver/pkg/line"
	"github.com/varas/numserver/pkg/report"
	"github.com/varas/numserver/pkg/repository"
	"github.com/varas/numserver/pkg/result"
)

// runtime services
type runtime struct {
	isUp          bool
	errHandle     errhandler.ErrHandler
	cancelRunners context.CancelFunc
	wgRunners     sync.WaitGroup
	listener      net.Listener
	connHandler   *connHandler
	resultWriter  *result.Writer
}

func (r *runtime) start(ctx context.Context, c config, errHandle errhandler.ErrHandler) (err error) {
	r.errHandle = errHandle

	ctx, r.cancelRunners = context.WithCancel(ctx)

	r.listener, err = net.Listen("tcp", fmt.Sprintf(":%d", c.port))
	if err != nil {
		return fmt.Errorf("cannot listen on socket tcp/%d error: %s", c.port, err.Error())
	}

	r.resultWriter, err = result.NewWriter(c.logPath, c.logFlushBatchSize)
	if err != nil {
		return fmt.Errorf("cannot create result writer: %s", err.Error())
	}

	numberRepository := repository.NewInMemoryRepository()
	reportRunner := report.NewRunner(c.reportFlushInterval)
	resultRunner := result.NewRunner(c.logFlushInterval, r.resultWriter, numberRepository)
	lineValidator, err := line.NewValidator()
	if err != nil {
		return fmt.Errorf("cannot create line validator: %s", err.Error())
	}

	r.connHandler = newConnHandler(r.errHandle, lineValidator, numberRepository, reportRunner)

	r.wgRunners = sync.WaitGroup{}
	r.wgRunners.Add(2)
	go func() {
		r.errHandle(reportRunner.Run(ctx))
		r.wgRunners.Done()
	}()
	go func() {
		r.errHandle(resultRunner.Run(ctx))
		r.wgRunners.Done()
	}()

	r.isUp = true

	return nil
}

func (r *runtime) stop() {
	if !r.isUp {
		return
	}
	r.isUp = false
	// stop server and runners
	r.errHandle(r.listener.Close())
	r.cancelRunners()
	// wait for runners to flush pending data
	r.wgRunners.Wait()
	r.errHandle(r.resultWriter.Close())
}
