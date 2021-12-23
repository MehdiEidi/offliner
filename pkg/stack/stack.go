// Package stack implements a stack data structure for strings. It uses a slice as its backing data structure.
package stack

import (
	"errors"
	"sync"
)

// Stack contains a slice of strings which is the backing data structure and also a mutex for keeping thread safety.
type Stack struct {
	Lock sync.Mutex
	Data []string
}

// New initializes and returns a new stack.
func New() *Stack {
	return &Stack{sync.Mutex{}, make([]string, 0, 50)}
}

// Push pushes a string onto the stack.
func (s *Stack) Push(str string) {
	s.Lock.Lock()
	defer s.Lock.Unlock()
	s.Data = append(s.Data, str)
}

// Pop removes the top element and returns it. If the stack is empty, it will return an error and an empty string.
func (s *Stack) Pop() (string, error) {
	s.Lock.Lock()
	defer s.Lock.Unlock()

	l := len(s.Data)

	if l == 0 {
		return "", errors.New("Pop() - stack is empty")
	}

	str := s.Data[l-1]
	s.Data = s.Data[:l-1]

	return str, nil
}

// Top returns the top element on the stack. It returns an empty string and an error if the stack is empty.
func (s *Stack) Top() (string, error) {
	s.Lock.Lock()
	defer s.Lock.Unlock()

	l := len(s.Data)

	if l == 0 {
		return "", errors.New("Top() - stack is empty")
	}

	return s.Data[l-1], nil
}

// Len returns the number of elements on the stack.
func (s *Stack) Len() int {
	return len(s.Data)
}
