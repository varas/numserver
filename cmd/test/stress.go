package main

import (
	"context"
	"fmt"
	"time"

	"log"
	"net"

	"flag"

	"sync"

	"github.com/varas/numserver/pkg/server"
)

var (
	port           = flag.Int("port", server.DefaultPort, fmt.Sprintf("-port %d", server.DefaultPort))
	file           = flag.String("file", server.DefaultLogFile, fmt.Sprintf("-file %s", server.DefaultLogFile))
	clientsAmount  = flag.Int("clients", 5, fmt.Sprintf("-clients 5"))
	testPeriodSecs = flag.Int("seconds", 10, fmt.Sprintf("-seconds 10"))
)

func main() {
	flag.Parse()

	numbersToWritePerClient := 6000000 // 6M
	numbersPerWrite := 1000
	testPeriod := time.Duration(*testPeriodSecs) * time.Second

	log.Printf("clients: %d, writes/client: %d, numbers/write: %d\n", *clientsAmount, numbersToWritePerClient, numbersPerWrite)

	srv := server.NewNumServer(*port, *file)

	log.Println("generating inputs...")

	inputs := generateInputs(numbersToWritePerClient, numbersPerWrite)

	log.Println("starting numserver...")

	go srv.Run(context.Background())

	<-srv.Ready

	log.Println("connecting clients...")

	testStart := make(chan struct{})
	clientsReady := sync.WaitGroup{}
	clientsReady.Add(*clientsAmount)

	for _, clientInput := range inputs {
		go func(clientInput [][]byte) {
			client, err := net.Dial("tcp", fmt.Sprintf(":%d", *port))
			handleErr(err)
			defer client.Close()

			// wait for test start
			clientsReady.Done()
			<-testStart

			for _, input := range clientInput {
				_, err = client.Write([]byte(input))
				handleErr(err)
			}
		}(clientInput)
	}

	clientsReady.Wait()
	close(testStart)

	log.Printf("start writing... (%.0f sec)", testPeriod.Seconds())

	timer := time.NewTimer(testPeriod)
	<-timer.C
	log.Println("time up! stopping...")
	close(srv.Stop)
	<-srv.Stopped
	log.Println("gracefully stopped")
}

// generates unique inputs
func generateInputs(numbersPerClient, numbersPerWrite int) (inputsPerClient map[int][][]byte) {
	var unique = 1
	var numsPerWrite string
	inputsPerClient = make(map[int][][]byte, *clientsAmount)

	for i := 0; i < *clientsAmount; i++ {
		inputsClient := [][]byte{}
		for ; unique <= numbersPerClient; unique++ {
			numsPerWrite = fmt.Sprintf("%s%09d\n", numsPerWrite, unique)

			if unique%numbersPerWrite == 0 || unique == numbersPerClient {
				inputsClient = append(inputsClient, []byte(numsPerWrite))
				numsPerWrite = ""
			}
		}
		inputsPerClient[i] = inputsClient
	}

	return
}

func handleErr(err error) {
	if err == nil {
		return
	}

	log.Fatalf("[error] %s", err.Error())
}
