package krongo

import (
	"fmt"
	"testing"
	"time"
)

func TestScheduler(t *testing.T) {
	sched := NewScheduler()

	j1 := Every(time.Second*3, func() error {
		fmt.Println("Hey :)")
		return nil
	})

	sched.AddJob(j1)
	sched.Start()
}
