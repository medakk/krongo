# krongo

A simple job scheduling library for go applications.

krongo lets you easily run tasks at specific times in your application, or create a task that repeatedly runs at periodic intervals. krongo has a simple API, available on [godoc](https://godoc.org/github.com/medakk/krongo).

krongo uses a ticker(that defaults to ticking every 500ms, but can be configured) that checks whether there are any events available to be run. This is useful for a system which requires batch-based jobs that need to be run occasionaly. For more fine-grained control of the timing, it is preferrable to use [time.Ticker](https://golang.org/pkg/time/#Ticker).

## Usage

#### Basic Usage
```
sched := krongo.NewScheduler()

// krongo.Every creates a repeating job
job1 := krongo.Every(1*time.Second, func() error {
	fmt.Println("Wimoweh")
	return nil
})
sched.AddJob(job1)

// krongo.At creates a job that runs once
job2 := krongo.At(time.Now().Add(10*time.Second), func() error {
	fmt.Println("In the jungle, the mighty jungle")
	fmt.Println("The lion sleeps tonight")
	return nil
})
sched.AddJob(job2)

// sched.Start blocks forever. You can run this in a goroutine,
// and continue to add jobs. Stop it by calling sched.Stop()
sched.Start()
```

#### Create custom error handlers
```
sched := krongo.NewScheduler()

// Set the function to run whenever an error occurs
sched.SetErrorHandler(func(err error) {
	log.Printf("failed to run job: %v", err)
})

// krongo.Every creates a repeating job
job := krongo.Every(time.Second, func() error {
	_, err := http.Get("http://www.example.com")
	if err != nil {
		return err
	}

	return nil
})
sched.AddJob(job)

sched.Start()
```

## Contributing

All issues and PRs are welcome.

## License

This project is licensed under the MIT license.