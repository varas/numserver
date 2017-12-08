package server

import (
	"context"

	"github.com/pkg/errors"
	"github.com/varas/numserver/pkg/errhandler"
)

// NumServer tcp server that store unique numbers writen
// It works as a bg daemon, so its API is based on channels to trigger graceful stop and wait for state completions
type NumServer struct {
	config    config
	runtime   *runtime
	errHandle errhandler.ErrHandler
	Ready     chan struct{} // enables to wait until ready
	Stop      chan struct{} // enables to gracefully stop the server
	Stopped   chan struct{} // enables to wait until stopped
}

// NewNumServer generates a new num-server
func NewNumServer(port int, logPath string) *NumServer {
	errHandle := errhandler.Logger("[error] ")

	return &NumServer{
		config:    *newConfig(port, logPath),
		runtime:   &runtime{}, // stateless runtime to enable restart
		errHandle: errHandle,
		Ready:     make(chan struct{}),
	}
}

// Run bootstraps the runtime so resilience could be added via recover, and runs the app
// context cancellation is aimed for fast teardown, for graceful stop use Stop channel instead
func (s *NumServer) Run(ctx context.Context) {
	s.Stop = make(chan struct{})
	s.Stopped = make(chan struct{})

	err := s.runtime.start(ctx, s.config, s.errHandle)
	if err != nil {
		s.errHandle(errors.Wrap(err, "error on start"))
		return
	}

	go s.waitForContextTermination(ctx)
	go s.waitForClientStop()

	close(s.Ready)

	<-s.runtime.stopped
	s.stop()
}

func (s *NumServer) stop() {
	if s.runtime.isUp.IsSet() {
		s.runtime.stop()
		close(s.Stopped)
	}
}

func (s *NumServer) waitForClientStop() {
	<-s.Stop
	s.stop()
}

func (s *NumServer) waitForContextTermination(ctx context.Context) {
	<-ctx.Done()
	s.stop()
}
