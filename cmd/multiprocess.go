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

// runMultiprocess starts scraping in a multiprocess fashion. It first scrapes the homepage filling the queue. It creates Unix sockets for each worker and runs the workers to scrape a single URL. Workers return extracted URLs via the Unix socket and they get collected.
func runMultiprocess(maxWorkers, maxpage int) {
	home, _ := urls.Dequeue()
	processURL(home)

	remainingPage := maxpage - 1

	if remainingPage < maxWorkers {
		maxWorkers = remainingPage
	}

	if urls.Len() < maxWorkers {
		maxWorkers = urls.Len()
	}

	for progress.Current() < maxpage {
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

		processes := make([]*exec.Cmd, maxWorkers)
		for i := 0; i < maxWorkers; i++ {
			link, _ := urls.Dequeue()

			connNum := strconv.Itoa(i)

			processes[i] = exec.Command("./process", link, baseDomain, connNum)
			processes[i].Stdout = os.Stdout
			processes[i].Start()
		}

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

func collectLinks(line string) {
	lineURLs := strings.Split(line, " ")
	for _, u := range lineURLs {
		if !visited.Has(u) {
			visited.Add(u)
			urls.Enqueue(u)
		}
	}
}
