// Package queue implements a queue data structure for strings. It uses slice as its backing data structure.
package queue

import (
	"errors"
	"sync"
)

// Queue contains a slice of strings which is the backing data structure and also a mutex for keeping thread safety.
type Queue struct {
	Lock sync.Mutex
	Data []string
}

// New initializes and returns a new stack.
func New() *Queue {
	return &Queue{sync.Mutex{}, make([]string, 0, 50)}
}

// Enqueue adds an element to the queue.
func (q *Queue) Enqueue(str string) {
	q.Lock.Lock()
	defer q.Lock.Unlock()
	q.Data = append(q.Data, str)
}

// Dequeue remove an element from the queue and returns it. If the stack is empty, it will return an error and an empty string.
func (q *Queue) Dequeue() (string, error) {
	q.Lock.Lock()
	defer q.Lock.Unlock()

	l := len(q.Data)

	if l == 0 {
		return "", errors.New("Dequeue() - queue is empty")
	}

	str := q.Data[0]
	q.Data = q.Data[1:]

	return str, nil
}

// First returns the first element on the queue. If the queue is empty, it will return an error and an empty string.
func (q *Queue) First() (string, error) {
	q.Lock.Lock()
	defer q.Lock.Unlock()

	if len(q.Data) == 0 {
		return "", errors.New("First() - queue is empty")
	}

	return q.Data[0], nil
}

// Len returns the number of elements on the queue.
func (q *Queue) Len() int {
	q.Lock.Lock()
	defer q.Lock.Unlock()
	return len(q.Data)
}
