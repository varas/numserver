package line

import (
	"fmt"
	"regexp"
)

// Validator validates number lines
type Validator struct {
	regex *regexp.Regexp
}

// NewValidator validates line is 9-digit or "termination" ending on carriage-return
func NewValidator() (*Validator, error) {
	regex, err := regexp.Compile(fmt.Sprintf(`^(\d{9}|%s)\n$`, terminationLine))

	return &Validator{
		regex: regex,
	}, err
}

// IsValidLine returns true if valid
func (r *Validator) IsValidLine(line string) bool {
	return r.regex.MatchString(line)
}
