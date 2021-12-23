package main

import "github.com/MehdiEidi/offliner/workerpool"

// chooseExecMethod chooses between serial, multithreaded, or multiprocessing depending on the flags.
func chooseExecMethod(serial, multiprocess bool, maxpage, maxWorkers int) {
	switch {
	case serial:
		for urls.Len() != 0 && progress.PageNum < maxpage {
			url, _ := urls.Pop()
			processURL(url, 0)
		}

	case multiprocess:
		runMultiprocess(maxWorkers, maxpage)

	default: // default is multithreading
		wp := workerpool.New(maxWorkers, processURL)
		go wp.Start()

		for progress.PageNum < maxpage {
			url, err := urls.Pop()
			if err != nil {
				continue
			}

			wp.AddTask(url)
		}
	}
}
