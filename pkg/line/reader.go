package line

import (
	"bufio"

	"fmt"

	"io"

	"strconv"
	"strings"

	"github.com/pkg/errors"
)

const terminationLine = "terminate"

// ErrTermination error returned on termination input read
var ErrTermination = errors.New("termination")

// Reader reads lines to return valid numbers or termination
type Reader struct {
	reader    bufio.Reader
	validator *Validator
}

// NewReader reads number lines
func NewReader(reader bufio.Reader, validator *Validator) *Reader {
	return &Reader{
		reader:    reader,
		validator: validator,
	}
}

// ReadNumberLine reads a valid line or returns error
// special returned errors:
// * io.EOF: on input end
// * terminationError: on termination input
func (r *Reader) ReadNumberLine() (number uint32, err error) {
	line, err := r.reader.ReadString('\n')
	if err == io.EOF {
		return
	}

	if err != nil {
		err = errors.Wrap(err, "cannot read input")
		return
	}

	if !r.validator.IsValidLine(line) {
		err = fmt.Errorf("invalid line: %s", line)
		return
	}

	line = strings.TrimSuffix(line, "\n")

	if line == terminationLine {
		err = ErrTermination
		return
	}

	numberInt, err := strconv.Atoi(line)
	if err != nil {
		err = errors.Wrapf(err, "cannot convert line to int %s", line)
		return
	}

	number = uint32(numberInt)

	return
}
