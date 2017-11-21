package line

import (
	"bufio"
	"strings"
	"testing"

	"io"

	"github.com/stretchr/testify/assert"
)

func TestLineReader_ReadsOneLineEachTime(t *testing.T) {
	theeLineInput := "123456789\nInvalidLine\n987654321\n"

	validator, err := NewValidator()
	assert.NoError(t, err)

	r := NewReader(*bufio.NewReader(strings.NewReader(theeLineInput)), validator)

	line, err := r.ReadNumberLine()

	assert.NoError(t, err)
	assert.Equal(t, uint32(123456789), line)

	_, err = r.ReadNumberLine()

	assert.Error(t, err)

	line, err = r.ReadNumberLine()

	assert.NoError(t, err)
	assert.Equal(t, uint32(987654321), line)
}

func TestLineReader_ReturnsEOFErrorOnEndOfFile(t *testing.T) {
	validator, err := NewValidator()
	assert.NoError(t, err)

	r := NewReader(*bufio.NewReader(strings.NewReader("")), validator)

	_, err = r.ReadNumberLine()

	assert.Equal(t, io.EOF, err)
}
