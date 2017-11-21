package server

import (
	"bufio"
	"context"
	"io"
	"net"

	"bitbucket.org/jhvaras/numserver/pkg/errhandler"
	"bitbucket.org/jhvaras/numserver/pkg/line"
	"bitbucket.org/jhvaras/numserver/pkg/report"
	"bitbucket.org/jhvaras/numserver/pkg/repository"
)

type connHandler struct {
	errHandle        errhandler.ErrHandler
	lineValidator    *line.Validator
	numberRepository repository.NumberRepository
	reportRunner     *report.Runner
}

func newConnHandler(
	errHandle errhandler.ErrHandler,
	lineValidator *line.Validator,
	numberRepo repository.NumberRepository,
	reportRunner *report.Runner,
) *connHandler {
	return &connHandler{
		errHandle:        errHandle,
		lineValidator:    lineValidator,
		numberRepository: numberRepo,
		reportRunner:     reportRunner,
	}
}

func (r *connHandler) run(ctx context.Context, conns <-chan net.Conn, termination chan struct{}) {
	for {
		select {
		case <-ctx.Done():
			return

		case <-termination:
			return

		case c := <-conns:
			r.handle(c, termination)
		}
	}
}

func (r *connHandler) handle(conn net.Conn, termination chan struct{}) {
	defer conn.Close()
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
