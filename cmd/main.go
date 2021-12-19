package main

import (
	"flag"
	"log"

	"github.com/MehdiEidi/offliner/pkg/set"
	"github.com/MehdiEidi/offliner/pkg/stack"
	"github.com/schollz/progressbar/v3"
)

var (
	// Stack of URLs to be processed.
	urls = stack.New()

	// Set of visited URLs.
	visited = set.New()

	// We only want to crawl the URLs of the same base domain.
	baseDomain string

	// Number of the pages processed. Should keep this number below maxPage.
	pageNum int

	// Progress bar for CLI
	bar *progressbar.ProgressBar
)

func main() {
	// Defining some flags
	homepage := flag.String("url", "", "URL of the homepage.")
	useProcesses := flag.Bool("a", false, "If the concurrency must be done using multi processing instead of multi threading.")
	serial := flag.Bool("s", false, "If the crawling must be done in non-concurrent fashion.")
	maxPage := flag.Int("n", 100, "Max number of pages to save.")
	maxWorkers := flag.Int("p", 50, "Max number of concurrent execution units.")
	flag.Parse()

	if *homepage == "" {
		log.Fatal("homepage URL cannot be empty.")
	}

	// Initializing the progress bar.
	bar = progressbar.Default(int64(*maxPage))

	if err := findBaseDomain(*homepage); err != nil {
		log.Fatal(err)
	}

	if err := createDir(); err != nil {
		log.Fatal(err)
	}

	urls.Push(*homepage)

	// Choose between serial, multithreaded, or multiprocessing depending on the flags.
	chooseRunningMethod(*serial, *useProcesses, *maxPage, *maxWorkers)
}
