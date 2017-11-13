package result

import (
	"fmt"
	"os"
)

// Writer writes results to file
type Writer struct {
	fd *os.File
}

// NewWriter ...
func NewWriter(filePath string) (*Writer, error) {
	output, err := os.Create(filePath)
	if err != nil {
		return nil, fmt.Errorf("cannot create log file: %s", err.Error())
	}

	return &Writer{
		fd: output,
	}, nil
}

// Write ...
func (r *Writer) Write(numbers []uint32) (err error) {
	lines := ""

	for _, n := range numbers {
		lines = fmt.Sprintf("%d%s%s", n, "\n", lines)
	}

	_, err = r.fd.WriteString(lines)

	return
}

// Close closes result file
func (r *Writer) Close() error {
	return r.fd.Close()
}
