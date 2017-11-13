package server

import (
	"context"

	"bitbucket.org/jhvaras/numserver/src/errhandler"
	"github.com/pkg/errors"
)

// NumServer ...
type NumServer struct {
	config    config
	runtime   runtime
	errHandle errhandler.ErrHandler
	Ready     chan struct{} // reads block until runtime bootstraps
}

// NewNumServer ...
func NewNumServer(port int, logPath string) *NumServer {
	return &NumServer{
		config:    *newConfig(port, logPath),
		errHandle: errhandler.Logger("[error] "),
		Ready:     make(chan struct{}),
	}
}

// Run bootstraps the runtime so resilience could be added via recover, and runs the app
func (s *NumServer) Run(ctx context.Context) {
	err := s.runtime.start(ctx, s.config, s.errHandle)
	if err != nil {
		s.errHandle(errors.Wrap(err, "error on start"))
	}
	defer s.runtime.stop()

	concurrency := newConcurrencyManager(s.config.concurrentClients)

	termination := make(chan struct{})
	go func() {
		<-termination
		s.runtime.stop()
		return
	}()

	close(s.Ready)

	for {
		conn, err := s.runtime.acceptConn()
		if err != nil {
			s.errHandle(errors.Wrap(err, "error accepting connections"))
			continue
		}

		concurrency.AddTaskOrWait()

		go func() {
			s.runtime.handle(ctx, conn, termination)

			concurrency.FinishTask()
		}()
	}
}
