package server

import (
	"context"
	"fmt"
	"math/rand"
	"net"
	"testing"
	"time"

	"os"

	"sync"
	"sync/atomic"

	"github.com/stretchr/testify/assert"
	"github.com/varas/numserver/pkg/errhandler"
)

var (
	testDataFolder = "testdata"
	testFilePath   = fmt.Sprintf("%s%s%s", testDataFolder, string(os.PathSeparator), DefaultLogFile)
	// inputs:
	validOneLineInput   = "314159265\n"
	validMultiLineInput = "007007009\n314159265\n"
	validInputs         = []string{validOneLineInput, validMultiLineInput}
	invalidInput        = "too large line\nshort\n"
	inputs              = append(validInputs, invalidInput)
)

func TestNumServer_SupportsClientWrites(t *testing.T) {
	for _, input := range validInputs {

		client, err := runServerAndClient(errhandler.Noop)
		if err != nil {
			t.Fatalf("cannot connect to server: %s", err.Error())
		}
		defer client.Close()

		n, err := client.Write([]byte(input))

		assert.NoError(t, err)
		assert.Equal(t, len(input), n)
	}
}

func TestNumServer_HandlesErrorsOnInvalidLines(t *testing.T) {
	wg := sync.WaitGroup{}
	wg.Add(2) // 2 invalid lines

	spyHandler, handledAmount := countHandler(&wg)

	client, err := runServerAndClient(spyHandler)
	if err != nil {
		t.Fatalf("cannot connect to server: %s", err.Error())
	}
	defer client.Close()

	_, err = client.Write([]byte(invalidInput))
	assert.NoError(t, err)

	wg.Wait()

	assert.Equal(t, int32(2), *handledAmount, "invalid input should cause an error being handled")
}

func TestNumServer_DoesNotHandleErrorsOnValidInput(t *testing.T) {
	wg := sync.WaitGroup{}

	spyHandler, handledAmount := countHandler(&wg)

	client, err := runServerAndClient(spyHandler)
	if err != nil {
		t.Fatalf("cannot connect to server: %s", err.Error())
	}
	defer client.Close()

	_, err = client.Write([]byte(validMultiLineInput))
	assert.NoError(t, err)

	wg.Wait()

	assert.Equal(t, int32(0), *handledAmount)
}

func TestNumServer_SupportsConcurrentClientWrites(t *testing.T) {
	numClients := DefaultConcurrentClients

	port := runServer(errhandler.Noop)

	wg := sync.WaitGroup{}
	wg.Add(numClients)

	raceStart := make(chan struct{})

	for i := 0; i < numClients; i++ {
		client, err := net.Dial("tcp", fmt.Sprintf(":%d", port))
		assert.NoError(t, err)
		defer client.Close()

		go func() {
			// wait for race start
			<-raceStart

			for _, input := range inputs {
				n, err := client.Write([]byte(input))

				assert.NoError(t, err)
				assert.Equal(t, len(input), n)
			}

			wg.Done()
		}()
	}

	close(raceStart)
	wg.Wait()
}

func runServer(errHandler errhandler.ErrHandler) (port int) {
	port = randPort()

	// should create file truncating if exists
	srv := NewNumServer(port, testFilePath)
	srv.errHandle = errHandler

	go srv.Run(context.Background())

	// wait for runtime start
	<-srv.Ready

	return
}

func runServerAndClient(errHandler errhandler.ErrHandler) (client net.Conn, err error) {
	port := runServer(errHandler)

	client, err = net.Dial("tcp", fmt.Sprintf(":%d", port))

	return
}

// range: 49152 to 65535
// https://www.iana.org/assignments/service-names-port-numbers/service-names-port-numbers.xhtml
func randPort() int {
	rand.Seed(int64(time.Now().Nanosecond()))

	return 49152 + rand.Intn(65535-49152)
}

func countHandler(wg *sync.WaitGroup) (errhandler.ErrHandler, *int32) {
	handledAmount := int32(0)

	return func(e error) {
		atomic.AddInt32(&handledAmount, 1)
		wg.Done()
	}, &handledAmount
}
