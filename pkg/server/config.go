package server

import "time"

// Default config values
const (
	DefaultPort                = 4000
	DefaultLogFile             = "numbers.log"
	DefaultLogFlushBatchSize   = 1000
	DefaultLogFlushInterval    = 1 * time.Second
	DefaultReportFlushInterval = 1 * time.Second
	DefaultConcurrentClients   = 5
)

type config struct {
	port    int
	logPath string
	// write numbers to file in batches
	logFlushBatchSize int
	// flush to log interval
	logFlushInterval time.Duration
	// report interval
	reportFlushInterval time.Duration
	// allowed concurrent clients
	concurrentClients int
}

func newConfig(port int, logPath string) *config {
	return &config{
		port:                port,
		logPath:             logPath,
		logFlushBatchSize:   DefaultLogFlushBatchSize,
		logFlushInterval:    DefaultLogFlushInterval,
		reportFlushInterval: DefaultReportFlushInterval,
		concurrentClients:   DefaultConcurrentClients,
	}
}
