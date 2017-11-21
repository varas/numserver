package main

import (
	"context"

	"log"

	"flag"
	"fmt"

	"github.com/varas/numserver/pkg/server"
)

var (
	port = flag.Int("port", server.DefaultPort, fmt.Sprintf("-port %d", server.DefaultPort))
	file = flag.String("file", server.DefaultLogFile, fmt.Sprintf("-file %s", server.DefaultLogFile))
	// we could also add other config params like:
	// * concurrentClients
	// * resultFlushInterval
	// * reportFlushInterval
)

func init() {
	flag.Parse()
}

func main() {
	srv := server.NewNumServer(*port, *file)

	// wait for runtime start
	go func() {
		<-srv.Ready
		log.Printf("numserver listening on tcp/%d and writing on %s", *port, *file)
	}()

	srv.Run(context.Background())
}
