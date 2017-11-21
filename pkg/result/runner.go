package result

import (
	"context"
	"time"

	"github.com/varas/numserver/pkg/repository"
)

// Runner prints a flush to standard output every 10 seconds
type Runner struct {
	interval   time.Duration
	writer     *Writer
	numberRepo repository.NumberRepository
}

// NewRunner creates a new daemon to write results on each interval
func NewRunner(interval time.Duration, writer *Writer, numberRepo repository.NumberRepository) *Runner {
	return &Runner{
		interval:   interval,
		writer:     writer,
		numberRepo: numberRepo,
	}
}

// Run runs writing results on each interval
func (r *Runner) Run(ctx context.Context) (err error) {
	ticker := time.NewTicker(r.interval)

	for {
		select {
		case <-ctx.Done():
			return r.flush()

		case <-ticker.C:
			err = r.flush()
			if err != nil {
				return
			}
		}
	}
}

func (r *Runner) flush() error {
	err := r.writer.Write(r.numberRepo.ExtractTransaction())
	if err != nil {
		r.numberRepo.Rollback()
		return err
	}

	r.numberRepo.Commit()

	return nil
}
