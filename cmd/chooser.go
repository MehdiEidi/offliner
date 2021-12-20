package main

import "github.com/MehdiEidi/offliner/workerpool"

// chooseRunningMethod chooses between serial, multithreaded, or multiprocessing depending on the flags.
func chooseRunningMethod(serial, useProcesses bool, maxPage, maxWorkers int) {
	switch {
	case serial:
		for urls.Len() != 0 && pageNum < maxPage {
			url, _ := urls.Pop()
			processURL(url, 0)
		}

	case useProcesses:
		runMultiProcess()

	default: // default is multithreading
		wp := workerpool.New(maxWorkers, processURL)
		go wp.Start()

		for pageNum < maxPage {
			url, err := urls.Pop()
			if err != nil {
				continue
			}

			wp.AddTask(url)
		}
	}
}
