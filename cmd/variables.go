package main

import (
	prog "github.com/MehdiEidi/offliner/pkg/progress"
	"github.com/MehdiEidi/offliner/pkg/queue"
	"github.com/MehdiEidi/offliner/pkg/set"
)

var (
	// Queue of the URLs to be processed.
	urls = queue.New()

	// Set of the already visited URLs.
	visited = set.New()

	// We only want to crawl the URLs of the same base domain.
	baseDomain string

	// Given command line flags.
	flags allFlags

	// Progress contains a progress bar and number of processed pages to keep track of progress.
	progress *prog.Progress
)

// allFlags contains given values for all the command line flags.
type allFlags struct {
	HomePage     string
	Maxworkers   int
	Maxpage      int
	Withfiles    bool
	Multiprocess bool
	Serial       bool
}
