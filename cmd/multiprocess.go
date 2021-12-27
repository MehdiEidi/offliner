package main

import (
	"log"
	"net"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"sync"
)

// runMultiprocess starts concurrent scraping in a multiprocess fashion. It first scrapes the homepage filling the queue. It creates Unix sockets for each worker and starts the workers to scrape a single URL. Workers return extracted URLs via the Unix socket in form of space separated string.
func runMultiprocess(maxWorkers, maxpage int) {
	home, _ := urls.Dequeue()
	processURL(home)

	remainingPage := maxpage - 1

	// In each iteration, maxWorkers number of pages get processed, so, if we have more workers than pages, we are going to end up processing more pages than desired.
	if remainingPage < maxWorkers {
		maxWorkers = remainingPage
	}

	if urls.Len() < maxWorkers {
		maxWorkers = urls.Len()
	}

	for progress.Current() < maxpage {
		// Creating sockets for each worker, the loop number is the ID of the socket.
		sockets := make([]net.Listener, maxWorkers)
		for i := 0; i < maxWorkers; i++ {
			connNum := strconv.Itoa(i)

			var err error
			sockets[i], err = net.Listen("unix", "/tmp/ipc"+connNum+".sock")
			if err != nil {
				log.Println(err)
			}
		}

		var wg sync.WaitGroup

		// Start accepting connections from the workers. Each worker will use the ID to connect to the appropriate socket.
		conns := make([]net.Conn, maxWorkers)
		for i := 0; i < maxWorkers; i++ {
			wg.Add(1)
			go func(connNum int) {
				wg.Done()

				var err error
				conns[connNum], err = sockets[connNum].Accept()
				if err != nil {
					log.Println(err)
				}
			}(i)
		}

		wg.Wait()

		// Starting the workers. For each worker, a link is taken from the queue and sent to it as a command line argument.
		processes := make([]*exec.Cmd, maxWorkers)
		for i := 0; i < maxWorkers; i++ {
			link, _ := urls.Dequeue()

			connNum := strconv.Itoa(i)

			processes[i] = exec.Command("./process", link, baseDomain, connNum)
			processes[i].Stdout = os.Stdout
			processes[i].Start()
		}

		// Iterating over the connections, reading, and collecting the data in which workers have been written. The data is a space separated line of URLs.
		for i := 0; i < maxWorkers; i++ {
			if conns[i] != nil {
				line := make([]byte, 1000000)

				n, err := conns[i].Read(line)
				if err != nil || n == 0 {
					continue
				}

				collectLinks(string(line[:n]))
			}
		}

		for i := 0; i < maxWorkers; i++ {
			processes[i].Wait()
		}

		for i := 0; i < maxWorkers; i++ {
			sockets[i].Close()
		}

		progress.Add(maxWorkers)

		remainingPage -= maxWorkers

		if remainingPage < maxWorkers {
			maxWorkers = remainingPage
		}

		if urls.Len() < maxWorkers {
			maxWorkers = urls.Len()
		}
	}
}

// collectLinks gets a line of space separated URLs, splits them into a slice, ranges over the slice, and adds the URLs to the queue if hasn't been visited before.
func collectLinks(line string) {
	lineURLs := strings.Split(line, " ")

	for _, u := range lineURLs {
		if !visited.Has(u) {
			visited.Add(u)
			urls.Enqueue(u)
		}
	}
}
