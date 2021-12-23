package main

import (
	"flag"
	"log"
	"net/url"

	prog "github.com/MehdiEidi/offliner/pkg/progress"
)

func main() {
	// Defining some flags
	homepage := flag.String("url", "", "URL of the homepage.")
	multiprocess := flag.Bool("a", false, "If the concurrency must be done using multi processing instead of multi threading.")
	serial := flag.Bool("s", false, "If the crawling must be done in non-concurrent fashion.")
	maxpage := flag.Int("n", 100, "Max number of pages to save.")
	maxWorkers := flag.Int("p", 50, "Max number of concurrent execution units.")
	flag.Parse()

	if *homepage == "" {
		log.Fatal("homepage URL cannot be empty.")
	}

	// Some initialization done.
	initialize(*homepage, *maxpage)

	// Choose between serial, multithreaded, or multiprocessing depending on the flags.
	chooseExecMethod(*serial, *multiprocess, *maxpage, *maxWorkers)
}

// initialize sets some global variables and initializes some stuff.
func initialize(homepage string, maxpage int) {
	u, err := url.Parse(homepage)
	if err != nil {
		log.Fatal(err)
	}

	// If homeURL has no scheme, default is http.
	if u.Scheme == "" {
		homepage = "http://" + homepage
	}

	if err := setHomeURL(homepage); err != nil {
		log.Fatal(err)
	}

	if err := setBaseDomain(homepage); err != nil {
		log.Fatal(err)
	}

	// Creating necessary directories.
	if err := createDirs(); err != nil {
		log.Fatal(err)
	}

	progress = prog.New(maxpage)

	urls.Push(homepage)
}
