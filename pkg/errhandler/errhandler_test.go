package errhandler

import (
	"bytes"
	"log"
	"os"
	"strings"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func TestErrLogger(t *testing.T) {
	var buf bytes.Buffer
	log.SetOutput(&buf)
	defer func() {
		log.SetOutput(os.Stderr)
	}()

	logger := Logger("")
	logger(errors.New("foo"))
	logger(errors.New("bar"))

	logged := buf.String()

	assert.True(t, strings.Contains(logged, "foo"))
	assert.True(t, strings.Contains(logged, "bar"))
	assert.Equal(t, 2, strings.Count(logged, "\n"), "should add carriage return on each error")

	// show log on verbose
	t.Log(logged)
}

func TestErrLogger_PrefixesLog(t *testing.T) {
	var buf bytes.Buffer
	log.SetOutput(&buf)
	defer func() {
		log.SetOutput(os.Stderr)
	}()

	Logger("[PREFIX] ")(errors.New("foo"))

	logged := buf.String()

	assert.True(t, strings.Contains(logged, "[PREFIX] "))

	// show log on verbose
	t.Log(logged)
}
