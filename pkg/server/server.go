package server

import (
	"context"

	"net"

	"os"
	"os/signal"
	"syscall"

	"bitbucket.org/jhvaras/numserver/pkg/errhandler"
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
		return
	}

	conns := make(chan net.Conn)
	termination := make(chan struct{})

	for w := s.config.concurrentClients; w >= 0; w-- {
		go s.runtime.connHandler.run(ctx, conns, termination)
	}

	go s.waitForClientTermination(termination)
	go s.waitForSystemTermination()

	close(s.Ready)

	for {
		conn, err := s.runtime.listener.Accept()
		if err != nil {
			if s.runtime.isUp {
				s.errHandle(errors.Wrap(err, "cannot accept connections"))
			}
			return
		}

		conns <- conn
	}
}

func (s *NumServer) waitForSystemTermination() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	<-c
	s.runtime.stop()
}

func (s *NumServer) waitForClientTermination(termination chan struct{}) {
	<-termination
	s.runtime.stop()
}
