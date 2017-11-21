package report

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReport_IncreaseUniqueTrue(t *testing.T) {
	r := Report{}

	r.Increase(true)

	assert.Equal(t, uint(1), r.uniqueDiff)
	assert.Equal(t, uint(1), r.uniqueTotal)
	assert.Equal(t, uint(0), r.duplicateDiff)
}

func TestReport_IncreaseUniqueFalse(t *testing.T) {
	r := Report{}

	r.Increase(false)

	assert.Equal(t, uint(0), r.uniqueDiff)
	assert.Equal(t, uint(0), r.uniqueTotal)
	assert.Equal(t, uint(1), r.duplicateDiff)
}

func TestReport_ReportTransactionCommit(t *testing.T) {
	r := Report{}

	r.Increase(true)

	_ = r.ReportTransaction()
	r.Commit()

	assert.Equal(t, uint(0), r.uniqueDiff)
	assert.Equal(t, uint(1), r.uniqueTotal)
	assert.Equal(t, uint(0), r.duplicateDiff)
}

func TestReport_ReportTransactionRollback(t *testing.T) {
	r := Report{}

	r.Increase(true)

	_ = r.ReportTransaction()
	r.Rollback()

	assert.Equal(t, uint(1), r.uniqueDiff)
	assert.Equal(t, uint(1), r.uniqueTotal)
	assert.Equal(t, uint(0), r.duplicateDiff)
}

func TestReport_SeveralPeriods(t *testing.T) {
	r := Report{}

	r.Increase(true)

	_ = r.ReportTransaction()
	r.Commit()

	r.Increase(true)

	assert.Equal(t, uint(1), r.uniqueDiff)
	assert.Equal(t, uint(2), r.uniqueTotal)
	assert.Equal(t, uint(0), r.duplicateDiff)
}
