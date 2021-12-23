package main

import (
	"bufio"
	"log"
	"os"
	"os/exec"
	"strings"
)

func runMultiprocess(maxWorkers int, maxpage int) {
	home, _ := urls.Pop()
	processURL(home, 0)

	temp, err := os.OpenFile("../temp/temp.txt", os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		log.Fatal(err)
	}

	tempMaxpage := maxpage

	// Because in each iteration, maxWorkers number of prosesses run concurrently, we cant have more workers than the number of maxpages.
	if maxpage < maxWorkers {
		maxWorkers = maxpage
	}

	for progress.Current() < maxpage {
		processes := make([]*exec.Cmd, maxWorkers)

		// Creating maxWorker number of concurrent processes. Each process scrapes a link.
		for i := 0; i < maxWorkers; i++ {
			link, err := urls.Pop()
			if err != nil {
				i--
				continue
			}

			processes[i] = exec.Command("./process", link, baseDomain)
			processes[i].Stdout = os.Stdout
			processes[i].Start()
		}

		// Waiting for all processes to finish.
		for i := 0; i < maxWorkers; i++ {
			processes[i].Wait()
		}

		// Read Lines of the processes output in temp file to a slice.
		var lines []string
		for i := 0; i < maxWorkers; i++ {
			scanner := bufio.NewScanner(temp)
			scanner.Scan()
			lines = append(lines, scanner.Text())
		}

		for i := 0; i < maxWorkers; i++ {
			lineURLs := strings.Split(lines[i], " ")
			for _, u := range lineURLs {
				if !visited.Has(u) {
					visited.Add(u)
					urls.Push(u)
				}
			}
		}

		// Reset the read/write offset of the temp file.
		temp.Seek(0, 0)

		progress.Add(maxWorkers)

		// Now we have saved maxWorker number of pages. We should update the maxWorkers number for the next iteration.
		tempMaxpage -= maxWorkers
		if tempMaxpage < maxWorkers {
			maxWorkers = tempMaxpage
		}
	}

	// Cleanup temp.
	if err = os.Remove("../temp/temp.txt"); err != nil {
		log.Fatal(err)
	}
}
