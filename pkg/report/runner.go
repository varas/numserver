package report

import (
	"context"
	"os"
	"time"
)

// Runner prints a report to standard output every 10 seconds
type Runner struct {
	interval time.Duration
	output   *os.File
	count    Report
}

// NewRunner creates a report runner daemon
func NewRunner(interval time.Duration) *Runner {
	return &Runner{
		interval: interval,
		output:   os.Stdout,
		count:    Report{},
	}
}

// Run runs reporting on each interval
func (r *Runner) Run(ctx context.Context) (err error) {
	ticker := time.NewTicker(r.interval)

	for {
		select {
		case <-ctx.Done():
			return r.report()

		case <-ticker.C:
			err = r.report()
			if err != nil {
				return
			}
		}
	}
}

// Increase increases report count
func (r *Runner) Increase(unique bool) {
	r.count.Increase(unique)
}

func (r *Runner) report() error {
	_, err := r.output.WriteString(r.count.ReportTransaction())
	if err != nil {
		r.count.Rollback()
		return err
	}
	r.count.Commit()

	return nil
}
