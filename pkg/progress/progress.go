// Package progress implements a thread-safe progress mechanism which will be used in crawler to keep track of the progress.
package progress

import (
	"sync"

	"github.com/schollz/progressbar/v3"
)

// Progress contains a mutex to ensure thread-safety, PageNum for keeping track of the number of the processed pages, and a progress bar for CLI.
type Progress struct {
	Lock    sync.Mutex
	PageNum int
	Bar     *progressbar.ProgressBar
}

// New initializes and returns a Progress. It uses default progressbar config.
func New(max int) *Progress {
	return &Progress{
		Lock:    sync.Mutex{},
		PageNum: 0,
		Bar:     progressbar.Default(int64(max)),
	}
}

// Add adds delta to the progress.
func (p *Progress) Add(delta int) {
	p.Lock.Lock()
	defer p.Lock.Unlock()

	p.PageNum++
	p.Bar.Add(delta)
}

// Current returns the current number of pages processed.
func (p *Progress) Current() int {
	p.Lock.Lock()
	defer p.Lock.Unlock()
	return p.PageNum
}
