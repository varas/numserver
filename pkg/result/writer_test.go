package result

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	testDataFolder = "testdata"
	testFilePath   = fmt.Sprintf("%s%s%s", testDataFolder, string(os.PathSeparator), "sample.log")
)

func TestWriter_Write(t *testing.T) {
	w, err := NewWriter(testFilePath)
	assert.NoError(t, err)

	assert.NoError(t, w.Write([]uint32{123456789, 001}))

	content, err := ioutil.ReadFile(testFilePath)
	assert.NoError(t, err)
	assert.Equal(t, "1\n123456789\n", string(content), "expected LIFO order with no leading zeroes")
}
