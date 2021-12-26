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

func runMultiprocess(maxWorkers, maxpage int) {
	home, _ := urls.Dequeue()
	processURL(home)
	progress.Add(1)

	remainingPage := maxpage - 1

	if remainingPage < maxWorkers {
		maxWorkers = remainingPage
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
				var err error
				conns[connNum], err = sockets[connNum].Accept()
				if err != nil {
					log.Println(err)
				}

				wg.Done()
			}(i)
		}

		processes := make([]*exec.Cmd, maxWorkers)
		for i := 0; i < maxWorkers; i++ {
			link, err := urls.Dequeue()
			if err != nil {
				i--
				continue
			}

			connNum := strconv.Itoa(i)

			processes[i] = exec.Command("./process", link, baseDomain, connNum)
			processes[i].Stdout = os.Stdout
			processes[i].Start()
		}

		wg.Wait()

		for i := 0; i < maxWorkers; i++ {
			if conns[i] != nil {
				var line []byte
				_, err := conns[i].Read(line)
				if err != nil {
					continue
				}

				collectLinks(string(line))
			} else {
				continue
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
	}

	defer func() {
		for i := 0; i < maxWorkers; i++ {
			os.Remove("/tmp/ipc" + strconv.Itoa(i) + ".sock")
		}
	}()
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
