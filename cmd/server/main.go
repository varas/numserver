package main

import (
	"context"

	"log"

	"bitbucket.org/jhvaras/numserver/src/server"
)

func main() {
	srv := server.NewNumServer(server.DefaultPort, server.DefaultFile)

	// wait for runtime start
	go func() {
		<-srv.Ready
		log.Printf("numserver listening tcp/%d and writing on %s", server.DefaultPort, server.DefaultFile)
	}()

	srv.Run(context.Background())

}
