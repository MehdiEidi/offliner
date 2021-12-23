package main

import (
	prog "github.com/MehdiEidi/offliner/pkg/progress"
	"github.com/MehdiEidi/offliner/pkg/set"
	"github.com/MehdiEidi/offliner/pkg/stack"
)

var (
	// Stack of URLs to be processed.
	urls = stack.New()

	// Set of already visited URLs.
	visited = set.New()

	// We only want to crawl the URLs of the same base domain.
	baseDomain string

	// URL of the homepage.
	homeURL string

	// Progress contains a progress bar and number of processed pages.
	progress *prog.Progress
)
