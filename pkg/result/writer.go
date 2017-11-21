package result

import (
	"fmt"
	"os"
)

// Writer writes results to file
type Writer struct {
	fd             *os.File
	flushBatchSize int
}

// NewWriter creates a new writer on the given file, writing bytes on batched size
func NewWriter(filePath string, flushBatchSize int) (*Writer, error) {
	output, err := os.Create(filePath)
	if err != nil {
		return nil, fmt.Errorf("cannot create log file: %s", err.Error())
	}

	return &Writer{
		fd:             output,
		flushBatchSize: flushBatchSize,
	}, nil
}

// Write writes numbers flushing to file on batches
func (r *Writer) Write(numbers []uint32) (err error) {
	lines := ""

	i := 0
	amount := len(numbers)
	for _, n := range numbers {
		lines = fmt.Sprintf("%d%s%s", n, "\n", lines)

		i++
		if i%r.flushBatchSize == 0 || i == amount {
			_, err = r.fd.WriteString(lines)
			if err != nil {
				return
			}
			lines = ""
		}
	}

	return
}

// Close closes result file
func (r *Writer) Close() error {
	return r.fd.Close()
}
