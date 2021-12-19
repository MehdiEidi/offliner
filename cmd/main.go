package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/MehdiEidi/offliner/pkg/set"
	"github.com/MehdiEidi/offliner/pkg/stack"
	"github.com/MehdiEidi/offliner/workerpool"
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

	if err := findBaseDomain(*homepage); err != nil {
		log.Fatal(err)
	}

	if err := createDir(); err != nil {
		log.Fatal(err)
	}

	// Process the homepage to initially fill the stack of URLs.
	err := process(*homepage, -1)
	if err != nil {
		log.Println(err)
	}

	// Choose between serial, multithreaded, or multiprocessing
	switch {
	case *serial:
		log.Println("--serial processing--")

		for urls.Len() != 0 {
			url, err := urls.Pop()
			if err != nil {
				log.Fatal(err)
			}

			err = process(url, -1)
			if err != nil {
				log.Println(err)
			}
		}

	case *useProcesses:
		// Todo

	default:
		fmt.Println("--multithreaded processing--")

		wp := workerpool.New(*maxWorkers, process)
		go wp.Start()

		for pageNum < *maxPage {
			url, err := urls.Pop()
			if err != nil {
				log.Println(err)
				break
			}
			wp.AddTask(url)
		}
	}
}
