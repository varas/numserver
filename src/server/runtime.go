package server

import (
	"context"
	"fmt"
	"net"

	"bufio"

	"io"

	"bitbucket.org/jhvaras/numserver/src/errhandler"
	"bitbucket.org/jhvaras/numserver/src/line"
	"bitbucket.org/jhvaras/numserver/src/report"
	"bitbucket.org/jhvaras/numserver/src/repository"
	"bitbucket.org/jhvaras/numserver/src/result"
)

// runtime services
type runtime struct {
	errHandle        errhandler.ErrHandler
	listener         net.Listener
	lineValidator    *line.Validator
	numberRepository repository.NumberRepository
	resultWriter     *result.Writer
	resultRunner     *result.Runner
	reportRunner     *report.Runner
	ctxCancel        context.CancelFunc
}

func (r *runtime) start(ctx context.Context, c config, errHandle errhandler.ErrHandler) error {
	r.errHandle = errHandle

	ctx, r.ctxCancel = context.WithCancel(ctx)
	var err error

	r.resultWriter, err = result.NewWriter(c.logPath)
	if err != nil {
		return fmt.Errorf("cannot create result writer: %s", err.Error())
	}

	r.numberRepository = repository.NewInMemoryRepository()

	r.lineValidator, err = line.NewValidator()
	if err != nil {
		return fmt.Errorf("cannot create line validator: %s", err.Error())
	}

	r.listener, err = net.Listen("tcp", fmt.Sprintf(":%d", c.port))
	if err != nil {
		return fmt.Errorf("cannot listen on socket tcp/%d error: %s", c.port, err.Error())
	}

	r.reportRunner = report.NewRunner(c.reportFlushInterval)

	go func() {
		r.errHandle(r.reportRunner.Run(ctx))
	}()

	r.resultRunner = result.NewRunner(c.resultFlushInterval, r.resultWriter, r.numberRepository)

	go func() {
		r.errHandle(r.resultRunner.Run(ctx))
	}()

	return nil
}

func (r *runtime) stop() {
	r.errHandle(r.listener.Close())
	r.ctxCancel()
	r.errHandle(r.resultWriter.Close())
}

func (r *runtime) acceptConn() (net.Conn, error) {
	return r.listener.Accept()
}

func (r *runtime) handle(ctx context.Context, conn net.Conn, termination chan struct{}) {
	reader := line.NewReader(*bufio.NewReader(conn), r.lineValidator)

	for {
		num, err := reader.ReadNumberLine()
		if err == io.EOF {
			return
		}

		if err == line.ErrTermination {
			close(termination)
			return
		}

		if err != nil {
			r.errHandle(err)
			continue
		}

		unique := r.numberRepository.AddNumber(num)
		r.reportRunner.Increase(unique)
	}
}
