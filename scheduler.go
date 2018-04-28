package krongo

import (
	"sync"
	"time"
)

type Scheduler struct {
	mu           *sync.Mutex
	jobs         []Job
	ticker       *time.Ticker
	errorHandler func(error)
	errorChannel chan error
	exitChannel  chan interface{}
}

const (
	defaultTickDuration = time.Second
)

func defaultErrorHandler(error) {}

// NewScheduler creates a new scheduler
func NewScheduler() *Scheduler {
	mu := &sync.Mutex{}
	jobs := make([]Job, 0)
	exitChannel := make(chan interface{}, 1)

	//TODO: too many errors might cause blocking on this channel
	errorChannel := make(chan error, 1024)

	s := Scheduler{
		mu:           mu,
		jobs:         jobs,
		ticker:       time.NewTicker(defaultTickDuration),
		errorHandler: defaultErrorHandler,
		errorChannel: errorChannel,
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
							sched.errorChannel <- err
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
// been started, and the current tick continues to get executed
func (sched *Scheduler) Stop() {
	sched.ticker.Stop()
	sched.exitChannel <- nil
}
