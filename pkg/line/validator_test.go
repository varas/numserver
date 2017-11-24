package line

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLineValidator_ReturnsTrueOnValidLines(t *testing.T) {
	linesWithNot10Len := []string{
		"123456789\n",
		"314159265\n",
		"007007009\n",
		"terminate\n",
	}

	v, err := NewValidator()
	assert.NoError(t, err)

	for _, l := range linesWithNot10Len {
		assert.True(t, v.IsValidLine(l))
	}
}

func TestLineValidator_ReturnsFalseWhenLineCharsNot10(t *testing.T) {
	linesWithNot10Len := []string{
		"\n",
		"1\n",
		"12345678\n",
		"1234567890\n",
	}

	v, err := NewValidator()
	assert.NoError(t, err)

	for _, l := range linesWithNot10Len {
		assert.False(t, v.IsValidLine(l))
	}
}

func TestLineValidator_ReturnsFalseOnNonDigitCharacters(t *testing.T) {
	linesWithNot10Len := []string{
		"a\n",
		"12345678a\n",
		"12345678A\n",
		"12345678ñ\n",
		"12345678-\n",
	}

	v, err := NewValidator()
	assert.NoError(t, err)

	for _, l := range linesWithNot10Len {
		assert.False(t, v.IsValidLine(l))
	}
}

func TestLineValidator_ReturnsFalseOnEOFWithoutCarriageReturn(t *testing.T) {
	v, err := NewValidator()
	assert.NoError(t, err)

	assert.False(t, v.IsValidLine("123456789"))
}
