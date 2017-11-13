package server

import "time"

// Default config values
const (
	DefaultPort                = 4000
	DefaultFile                = "numbers.log"
	DefaultConcurrentClients   = 5
	DefaultResultFlushInterval = 1 * time.Second
	DefaultReportFlushInterval = 10 * time.Second
)

type config struct {
	port                int
	logPath             string
	concurrentClients   int
	resultFlushInterval time.Duration
	reportFlushInterval time.Duration
}

func newConfig(port int, logPath string) *config {
	return &config{
		port:                port,
		logPath:             logPath,
		concurrentClients:   DefaultConcurrentClients,
		resultFlushInterval: DefaultResultFlushInterval,
		reportFlushInterval: DefaultReportFlushInterval,
	}
}
