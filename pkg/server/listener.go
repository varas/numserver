package server

import (
	"context"

	"fmt"
	"net"

	"github.com/pkg/errors"
)

// Listener listens for connections and sends to output channel
type Listener struct {
	listener net.Listener
	conns    chan<- net.Conn
}

// NewListener creates new connection listener on given tcp port
func NewListener(port int, conns chan<- net.Conn) (*Listener, error) {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return nil, errors.Wrapf(err, "cannot listen on socket tcp/%d", port)
	}

	return &Listener{
		listener: listener,
		conns:    conns,
	}, nil
}

// Listen listens for new connections
func (s *Listener) Listen(ctx context.Context) error {
	go s.waitForContextTermination(ctx)

	for {
		conn, err := s.listener.Accept()
		if err != nil {
			return err
		}

		s.conns <- conn
	}
}

// Stop stops the listener gracefully
func (s *Listener) Stop() {
	_ = s.listener.Close()
	return
}

func (s *Listener) waitForContextTermination(ctx context.Context) {
	<-ctx.Done()
	s.Stop()
}
