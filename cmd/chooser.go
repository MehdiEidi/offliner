package main

import "github.com/MehdiEidi/offliner/pkg/workerpool"

// chooseExecMethod chooses between serial, multithreaded, or multiprocessing depending on the flags.
func chooseExecMethod() {
	switch {
	case flags.Serial:
		for urls.Len() != 0 && progress.Current() < flags.Maxpage {
			url, _ := urls.Dequeue()
			processURL(url)
		}

	case flags.Multiprocess:
		runMultiprocess(flags.Maxworkers, flags.Maxpage)

	default: // default is multithreading
		wp := workerpool.New(flags.Maxworkers, processURL)
		go wp.Start()

		for progress.Current() < flags.Maxpage {
			url, err := urls.Dequeue()
			if err != nil {
				continue
			}

			wp.AddTask(url)
		}
	}
}
