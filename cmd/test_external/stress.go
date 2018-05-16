package main

import (
"flag"
"fmt"
"log"
"net"
"os"
"runtime/pprof"
"sync"
"time"














"github.com/varas/numserver/pkg/server"
)

var (
	port           = flag.Int("port", server.DefaultPort, fmt.Sprintf("-port %d", server.DefaultPort))
	clientsAmount  = flag.Int("clients", 5, fmt.Sprintf("-clients 5"))
	testPeriodSecs = flag.Int("seconds", 10, fmt.Sprintf("-seconds 10"))
	cpuProfile     = flag.String("cpuprofile", "", "write cpu profile to file, -cpuprofile stress.cpu")
)

func main() {
	flag.Parse()

	numbersToWritePerClient := 6000000 // 6M
	numbersPerWrite := 1000
	testPeriod := time.Duration(*testPeriodSecs) * time.Second

	log.Printf("clients: %d, writes/client: %d, numbers/write: %d\n", *clientsAmount, numbersToWritePerClient, numbersPerWrite)


	log.Println("generating inputs...")

	inputs := generateInputs(numbersToWritePerClient, numbersPerWrite)

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

	if *cpuProfile != "" {
		f, err := os.Create(*cpuProfile)
		if err != nil {
			log.Fatal("could not create CPU profile: ", err)
		}
		if err := pprof.StartCPUProfile(f); err != nil {
			log.Fatal("could not start CPU profile: ", err)
		}
		defer pprof.StopCPUProfile()
	}

	clientsReady.Wait()
	close(testStart)

	log.Printf("start writing... (%.0f sec)", testPeriod.Seconds())

	timer := time.NewTimer(testPeriod)
	<-timer.C
	log.Println("time up! stopping...")
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
