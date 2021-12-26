// Package workerpool implements a workerpool of specific number of workers to apply a single task on a queue of URLs which are passed using a channel of string.
package workerpool

// Workerpool contains the max number of workers, a channel of strings acting as the task queue which will contain URLs, and a function which acts as the worker that will be applied on the URLs.
type Workerpool struct {
	MaxWorkers int
	URLQueue   chan string
	Worker     func(string) error
}

// New initializes and returns a workerpool.
func New(maxWorkers int, worker func(string) error) *Workerpool {
	return &Workerpool{
		MaxWorkers: maxWorkers,
		URLQueue:   make(chan string),
		Worker:     worker,
	}
}

// AddTask enqueues a new task (URL) on the task channel (URLQueue).
func (wp *Workerpool) AddTask(task string) {
	wp.URLQueue <- task
}

// Start creates MaxWorker number of goroutines. Each goroutine ranges over the task channel, fetches tasks (URLs), and applies the Worker() on it.
func (wp *Workerpool) Start() {
	for i := 0; i < wp.MaxWorkers; i++ {
		go func() {
			for url := range wp.URLQueue {
				wp.Worker(url)
			}
		}()
	}
}
