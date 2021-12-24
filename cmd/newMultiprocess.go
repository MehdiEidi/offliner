package main

import (
	"log"
	"net"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func runMultiprocess2(maxWorkers, maxpage int) {
	home, _ := urls.Pop()
	processURL(home, 0)

	tempMaxpage := maxpage

	if maxpage < maxWorkers {
		maxWorkers = maxpage - 1
	}

	for progress.Current() < maxpage {
		if maxWorkers < 0 {
			maxWorkers = tempMaxpage
		}

		sockets := make([]net.Listener, maxWorkers)
		for i := 0; i < maxWorkers; i++ {
			colNum := strconv.Itoa(i)
			var err error
			sockets[i], err = net.Listen("unix", "/tmp/ipc"+colNum+".sock")
			if err != nil {
				log.Println(err)
			}
		}

		conns := make([]net.Conn, maxWorkers)
		for i := 0; i < maxWorkers; i++ {
			go func(connNum int) {
				var err error
				conns[connNum], err = sockets[connNum].Accept()
				if err != nil {
					log.Println(err)
				}
			}(i)
		}

		processes := make([]*exec.Cmd, maxWorkers)
		for i := 0; i < maxWorkers; i++ {
			link, _ := urls.Pop()

			connNum := strconv.Itoa(i)

			processes[i] = exec.Command("./process2", link, baseDomain, connNum)
			processes[i].Stdout = os.Stdout
			processes[i].Start()
		}

		for i := 0; i < maxWorkers; i++ {
			var line []byte
			if conns[i] == nil {
				i--
				continue
			}
			conns[i].Read(line)

			collectLinks(string(line))
		}

		for i := 0; i < maxWorkers; i++ {
			processes[i].Wait()
		}

		for i := 0; i < maxWorkers; i++ {
			sockets[i].Close()
			conns[i].Close()
		}

		progress.Add(maxWorkers)

		tempMaxpage -= maxWorkers
		if tempMaxpage < maxWorkers {
			maxWorkers = tempMaxpage - 1
		}
	}
}

func collectLinks(line string) {
	lineURLs := strings.Split(string(line), " ")
	for _, u := range lineURLs {
		if !visited.Has(u) {
			visited.Add(u)
			urls.Push(u)
		}
	}
}
