package report

import (
	"fmt"
	"sync"
)

// Report stores the counts to be reported, supporting concurrency
// * The difference since the last report of the count of new unique numbers that have been received.
// * The difference since the last report of the count of new duplicate numbers that have been received.
// * The total number of unique numbers received for this run of the Application.
// * Example text: Received 50 unique numbers, 2 duplicates. Unique total: 567231
type Report struct {
	uniqueDiff    uint
	duplicateDiff uint
	uniqueTotal   uint
	sync.Mutex
}

// Increase ...
func (r *Report) Increase(unique bool) {
	r.Lock()
	if unique {
		r.uniqueDiff++
		r.uniqueTotal++
	} else {
		r.duplicateDiff++
	}
	r.Unlock()
}

// ReportTransaction retrieves report as human readable text starting a transaction to be committed or rollbacked
func (r *Report) ReportTransaction() string {
	r.Lock()
	return fmt.Sprintf("Received %d unique numbers, %d duplicates. Unique total: %d\n",
		r.uniqueDiff,
		r.duplicateDiff,
		r.uniqueTotal,
	)
}

// Commit unlocks and reset a new period count
func (r *Report) Commit() {
	r.uniqueDiff = 0
	r.duplicateDiff = 0

	r.Unlock()
}

// Rollback unlocks without data removal
func (r *Report) Rollback() {
	r.Unlock()
}
