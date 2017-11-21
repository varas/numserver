package result

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"strings"

	"github.com/stretchr/testify/assert"
)

var (
	testDataFolder = "testdata"
	testFilePath   = fmt.Sprintf("%s%s%s", testDataFolder, string(os.PathSeparator), "sample.log")
)

func TestWriter_Write(t *testing.T) {
	numbers := []uint32{123456789, 001, 002}

	w, err := NewWriter(testFilePath, 2)
	assert.NoError(t, err)

	assert.NoError(t, w.Write(numbers))

	assertFileContains(t, testFilePath, numbers)
}

func assertFileContains(t *testing.T, filePath string, expectedNumbers []uint32) {
	content, err := ioutil.ReadFile(filePath)
	assert.NoError(t, err)

	linesAmount := strings.Count(string(content), "\n")
	assert.Equal(t, len(expectedNumbers), linesAmount, "lines amount does not match numbers amount")

	for _, n := range expectedNumbers {
		assert.True(t, strings.Contains(string(content), fmt.Sprintf("%d\n", n)))
	}
}
