package main

import (
	"context"
	"fmt"
	"time"

	"log"
	"net"

	"flag"

	"github.com/varas/numserver/pkg/server"
)

var (
	port           = flag.Int("port", server.DefaultPort, fmt.Sprintf("-port %d", server.DefaultPort))
	file           = flag.String("file", server.DefaultLogFile, fmt.Sprintf("-file %s", server.DefaultLogFile))
	clientsAmount  = flag.Int("clients", 1, fmt.Sprintf("-clients 1"))
	testPeriodSecs = flag.Int("seconds", 10, fmt.Sprintf("-seconds 10"))
)

func main() {
	flag.Parse()

	numbersToWrite := 10000000 // 10M
	numbersPerWrite := 1000
	testPeriod := time.Duration(*testPeriodSecs) * time.Second

	srv := server.NewNumServer(*port, *file)

	log.Println("generating inputs...")

	inputs := generateInputs(numbersToWrite, numbersPerWrite)

	log.Println("starting numserver...")

	ctx, cancel := context.WithCancel(context.Background())
	go srv.Run(ctx)

	<-srv.Ready

	log.Println("connecting clients...")

	testStart := make(chan struct{})

	for i := 0; i < *clientsAmount; i++ {
		go func() {
			client, err := net.Dial("tcp", fmt.Sprintf(":%d", *port))
			handleErr(err)
			defer client.Close()

			// wait for test start
			<-testStart

			for _, input := range inputs {
				_, err = client.Write([]byte(input))
				handleErr(err)
			}
		}()
	}

	close(testStart)
	timer := time.NewTimer(testPeriod)

	log.Printf("start writing... (%f sec)", testPeriod.Seconds())

	<-timer.C
	cancel()
}

// unique number reference to generate unique inputs
var num = 1

// generates unique inputs as it is the worst case scenario
func generateInputs(numbersAmount, numbersPerWrite int) (inputs [][]byte) {
	var numbers string

	for ; num <= numbersAmount; num++ {
		numbers = fmt.Sprintf("%s%09d\n", numbers, num)

		if num%numbersPerWrite == 0 || num == numbersAmount {
			inputs = append(inputs, []byte(numbers))
			numbers = ""
		}
	}

	return
}

func handleErr(err error) {
	if err == nil {
		return
	}

	log.Fatalf("[error] %s", err.Error())
}
