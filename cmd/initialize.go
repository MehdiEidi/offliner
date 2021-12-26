package main

import (
	"log"
	"net/url"
	"os"

	prog "github.com/MehdiEidi/offliner/pkg/progress"
)

// initialize sets some global variables and initializes some stuff.
func initialize(homepage string, maxworkers, maxpages int, multiprocess, serial, withfiles bool) {
	flags = allFlags{
		HomePage:     homepage,
		Maxworkers:   maxworkers,
		Maxpage:      maxpages,
		Multiprocess: multiprocess,
		Serial:       serial,
		Withfiles:    withfiles,
	}

	u, err := url.Parse(homepage)
	if err != nil {
		log.Fatal(err)
	}

	if u.Scheme == "" {
		log.Fatal("given URL doesn't have any scheme.")
	}

	baseDomain = u.Host
	rootURL = u.Scheme + "://" + u.Host

	// Creating necessary directories.
	if err := createDirs(); err != nil {
		log.Fatal(err)
	}

	progress = prog.New(maxpages)

	urls.Enqueue(homepage)
}

// createDirs contains a slice of the directories to be created, it ranges over the slice and creates those directories.
func createDirs() error {
	// Slice of the directories to be created. Order is important.
	dirs := []string{"../output/", "../output/" + baseDomain, "../output/" + baseDomain + "/pages/", "../output/" + baseDomain + "/static/", "../output/" + baseDomain + "/static/css/", "../output/" + baseDomain + "/static/js/", "../output/" + baseDomain + "/static/files/", "../output/" + baseDomain + "/static/img/"}

	for _, d := range dirs {
		if _, err := os.Stat(d); os.IsNotExist(err) {
			if err = os.Mkdir(d, os.ModePerm); err != nil {
				return err
			}
		}
	}

	return nil
}
