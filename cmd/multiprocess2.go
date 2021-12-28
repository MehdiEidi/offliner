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

func runMultiprocess2(maxWorkers, maxpage int) {
	home, _ := urls.Dequeue()
	processURL(home)

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
		connNum := strconv.Itoa(i)

		processes[i] = exec.Command("./process2", connNum, baseDomain)
		processes[i].Stdout = os.Stdout
		processes[i].Start()
	}

	wg.Wait()

	for progress.Current() < maxpage {
		written := 0
		for written < maxWorkers && urls.Len() != 0 {
			link, _ := urls.Dequeue()
			conns[written].Write([]byte(link))
			written++
		}

		for i := 0; i < written-1; i++ {
			if conns[i] != nil {
				line := make([]byte, 1000000)

				n, err := conns[i].Read(line)
				progress.Add(1)
				if err != nil || n == 0 {
					continue
				}

				collectLinks2(string(line[:n]))
			}
		}
	}

	for i := 0; i < maxWorkers; i++ {
		conns[i].Write([]byte("done"))
	}

	for i := 0; i < maxWorkers; i++ {
		processes[i].Wait()
	}

	for i := 0; i < maxWorkers; i++ {
		sockets[i].Close()
	}
}

func collectLinks2(line string) {
	lineURLs := strings.Split(line, " ")

	for _, u := range lineURLs {
		if !visited.Has(u) {
			visited.Add(u)
			urls.Enqueue(u)
		}
	}
}
