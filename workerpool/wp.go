// Package workerpool just implements a workerpool.
package workerpool

type Workerpool struct {
	MaxWorkers int
	URLQueue   chan string
	Task       func(string, int) error
}

func (wp *Workerpool) AddTask(task string) {
	wp.URLQueue <- task
}

func (wp *Workerpool) Start() {
	for i := 0; i < wp.MaxWorkers; i++ {
		wID := i + 1

		go func(workerID int) {
			for url := range wp.URLQueue {
				wp.Task(url, workerID)
			}
		}(wID)
	}
}

func New(maxWorkers int, task func(string, int) error) Workerpool {
	return Workerpool{
		MaxWorkers: maxWorkers,
		URLQueue:   make(chan string),
		Task:       task,
	}
}
