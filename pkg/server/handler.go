package server

import (
	"bufio"
	"context"
	"io"
	"net"

	"github.com/varas/numserver/pkg/errhandler"
	"github.com/varas/numserver/pkg/line"
	"github.com/varas/numserver/pkg/report"
	"github.com/varas/numserver/pkg/repository"
)

type connHandler struct {
	errHandle        errhandler.ErrHandler
	lineValidator    *line.Validator
	numberRepository repository.NumberRepository
	report           *report.Report
	conns            <-chan net.Conn
	terminate        chan struct{}
}

func newConnHandler(
	errHandle errhandler.ErrHandler,
	lineValidator *line.Validator,
	numberRepo repository.NumberRepository,
	report *report.Report,
	conns <-chan net.Conn,
	terminate chan struct{},
) *connHandler {
	return &connHandler{
		errHandle:        errHandle,
		lineValidator:    lineValidator,
		numberRepository: numberRepo,
		report:           report,
		conns:            conns,
		terminate:        terminate,
	}
}

func (r *connHandler) run(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return

		case c, open := <-r.conns:
			if !open {
				return
			}
			r.handle(ctx, c)
		}
	}
}

// context unhandled here to avoid data loss, as client has no guarantees of sent data is processed on service stop
func (r *connHandler) handle(ctx context.Context, conn net.Conn) {
	defer conn.Close()
	reader := line.NewReader(*bufio.NewReader(conn), r.lineValidator)

	for {
		num, err := reader.ReadNumberLine()
		if err == io.EOF {
			return
		}

		if err == line.ErrTermination {
			close(r.terminate)
			return
		}

		if err != nil {
			r.errHandle(err)
			continue
		}

		unique := r.numberRepository.AddNumber(num)
		r.report.Increase(unique)
	}
}
