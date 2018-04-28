package krongo

import (
	"sync"
	"time"
)

// Scheduler handles the scheduling of different jobs
type Scheduler struct {
	mu           *sync.Mutex
	jobs         []Job
	ticker       *time.Ticker
	errorHandler func(error)
	exitChannel  chan interface{}
}

const (
	defaultTickDuration = 500 * time.Millisecond
)

func defaultErrorHandler(error) {}

// NewScheduler creates a new scheduler
func NewScheduler() *Scheduler {
	mu := &sync.Mutex{}
	jobs := make([]Job, 0)
	exitChannel := make(chan interface{}, 1)

	s := Scheduler{
		mu:           mu,
		jobs:         jobs,
		ticker:       time.NewTicker(defaultTickDuration),
		errorHandler: defaultErrorHandler,
		exitChannel:  exitChannel,
	}

	return &s
}

// AddJob adds a new job to the scheduler
func (sched *Scheduler) AddJob(j Job) {
	sched.mu.Lock()
	defer sched.mu.Unlock()

	sched.jobs = append(sched.jobs, j)
}

// Start begins the scheduler. This is an indefinitely blocking function.
// All jobs are processed in a separated goroutine, so if any repeating job
// runs longer than the tick duration, it will be run again at the next tick
func (sched *Scheduler) Start() {
	for {
		select {
		case <-sched.ticker.C:
			sched.mu.Lock()

			i := 0
			for i < len(sched.jobs) {
				job := sched.jobs[i]
				now := time.Now()
				if job.ShouldRun(now) {

					go func(t time.Time, j Job) {
						err := j.Run(t)
						if err != nil {
							sched.errorHandler(err)
						}
					}(now, job)

					if job.DeleteOnRun() {
						sched.jobs[i] = sched.jobs[len(sched.jobs)-1]
						sched.jobs = sched.jobs[:len(sched.jobs)-1]

						// skip to next iteration, since i already indexes
						// the next job
						continue
					}
				}
				i++
			}

			sched.mu.Unlock()

		case <-sched.exitChannel:
			break
		}
	}
}

// Stop stops the scheduler. This does not cancel any jobs that have already
// been started. If a tick is already in progress, the scheduler is stopped
// at the end of the tick
func (sched *Scheduler) Stop() {
	sched.ticker.Stop()
	sched.exitChannel <- nil
}

// SetTickerDuration is used to set the minimum duration between which to
// check of jobs and run them. The default is 500ms
func (sched *Scheduler) SetTickerDuration(tickDuration time.Duration) {
	sched.mu.Lock()
	defer sched.mu.Unlock()

	sched.ticker = time.NewTicker(tickDuration)
}

// SetErrorHandler is used to set the function that handlers errors from jobs
func (sched *Scheduler) SetErrorHandler(f func(error)) {
	sched.mu.Lock()
	defer sched.mu.Unlock()

	sched.errorHandler = f
}
