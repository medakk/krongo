package krongo

import "time"

// Job is an interface that should be satisfied by all items
// pushed to the scheduler
type Job interface {
	// Given the time, return whether or not this job should run
	ShouldRun(now time.Time) bool

	// Run the job, and return a bool signifying whether or not this
	// job should be removed from the scheduler, and an error if the
	// job that was run failed.
	Run(now time.Time) error

	// DeleteOnRun returns whether or not the job is considered finished
	// after calling run
	DeleteOnRun() bool
}

// Every creates a new job that runs `f` every `t`
func Every(t time.Duration, f func() error) Job {
	j := repeatedJob{
		duration: t,
		lastRun:  time.Unix(0, 0),
		f:        f,
	}
	return &j
}

// At creates a new job that runs once at time `t`.
// If the time t has already passed, the job will choose to run at the
// next tick
func At(t time.Time, f func() error) Job {
	j := oneShotJob{
		when: t,
		f:    f,
	}
	return &j
}

type repeatedJob struct {
	duration time.Duration
	lastRun  time.Time
	f        func() error
}

func (j *repeatedJob) ShouldRun(now time.Time) bool {
	return now.Sub(j.lastRun) >= j.duration
}

func (j *repeatedJob) Run(now time.Time) error {
	j.lastRun = now
	return j.f()
}

func (j *repeatedJob) DeleteOnRun() bool {
	return false
}

type oneShotJob struct {
	when time.Time
	f    func() error
}

func (j *oneShotJob) ShouldRun(now time.Time) bool {
	return now.After(j.when) || now.Equal(j.when)
}

func (j *oneShotJob) Run(now time.Time) error {
	return j.f()
}

func (j *oneShotJob) DeleteOnRun() bool {
	return true
}
