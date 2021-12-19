// Package workerpool just implements a workerpool.
package workerpool

// Workerpool contains the max number of workers, a channel of URL acting as the task queue, and a function which acts as the worker.
type Workerpool struct {
	MaxWorkers int
	URLQueue   chan string
	Worker     func(string, int) error
}

// New initializes and returns a workerpool.
func New(maxWorkers int, worker func(string, int) error) *Workerpool {
	return &Workerpool{
		MaxWorkers: maxWorkers,
		URLQueue:   make(chan string),
		Worker:     worker,
	}
}

// AddTask enqueues a new task on the task channel.
func (wp *Workerpool) AddTask(task string) {
	wp.URLQueue <- task
}

// Start spins maxWorker number of goroutines. Each goroutine ranges over the task channel, fetches tasks, and runs it.
func (wp *Workerpool) Start() {
	for i := 0; i < wp.MaxWorkers; i++ {
		wID := i + 1

		go func(workerID int) {
			for url := range wp.URLQueue {
				wp.Worker(url, workerID)
			}
		}(wID)
	}
}
