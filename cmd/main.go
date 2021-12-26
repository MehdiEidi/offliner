package main

import (
	"flag"
	"log"
)

func main() {
	// Defining some flags.
	homepage := flag.String("url", "", "URL of the homepage.")
	multiprocess := flag.Bool("a", false, "If the concurrency must be done using multi processing instead of multi threading.")
	serial := flag.Bool("s", false, "If the crawling must be done in a non-concurrent fashion.")
	withfiles := flag.Bool("f", false, "If the static files must be downloaded too.")
	maxpages := flag.Int("n", 100, "Max number of pages to be saved.")
	maxworkers := flag.Int("p", 50, "Max number of concurrent execution units.")
	flag.Parse()

	if *homepage == "" {
		log.Fatal("you must provide a URL to start the scraping.")
	}

	initialize(*homepage, *maxworkers, *maxpages, *multiprocess, *serial, *withfiles)

	// Choose between serial, multithreaded, or multiprocessing depending on the flags.
	chooseExecMethod()
}
