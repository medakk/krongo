package krongo

import (
	"errors"
	"testing"
	"time"
)

//TODO: Create a method to mock the time values in the scheduler, so
//that it is possible to test the actual scheduling

func TestErrorHandler(t *testing.T) {
	sched := NewScheduler()

	ok := false

	sched.SetErrorHandler(func(err error) {
		ok = true
	})

	job := At(time.Now(), func() error {
		return errors.New("")
	})

	sched.AddJob(job)
	sched.SetTickerDuration(100 * time.Millisecond)
	go sched.Start()

	time.Sleep(200 * time.Millisecond)
	sched.Stop()

	if !ok {
		t.Errorf("error handler was not invoked")
	}
}
