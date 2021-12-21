package main

import (
	"bufio"
	"log"
	"os"
	"os/exec"
	"strings"
)

func runMultiProcess(maxWorkers int, maxPage int) {
	home, _ := urls.Pop()

	processURL(home, 0)

	temp, err := os.OpenFile("../temp/temp.txt", os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		log.Fatal(err)
	}

	tempMaxPage := maxPage

	if maxPage < maxWorkers {
		if maxPage > 30 {
			maxWorkers = 30
		} else {
			maxWorkers = maxPage - 1
		}
	}

	for pageNum < maxPage {
		var urlLines []string
		cmds := make([]*exec.Cmd, maxWorkers)

		for i := 0; i < maxWorkers; i++ {
			link, _ := urls.Pop()
			cmds[i] = exec.Command("./process", link, baseDomain)
			cmds[i].Start()
		}

		for i := 0; i < maxWorkers; i++ {
			cmds[i].Wait()
		}

		for i := 0; i < maxWorkers; i++ {
			scanner := bufio.NewScanner(temp)
			scanner.Scan()
			urlLines = append(urlLines, scanner.Text())
		}

		var u []string

		for i := 0; i < maxWorkers; i++ {
			t := strings.Split(urlLines[i], " ")
			for j := 0; j < len(t); j++ {
				u = append(u, t[j])
			}
		}

		for i := 0; i < len(u); i++ {
			l := u[i]
			if !visited.Has(l) {
				urls.Push(l)
				visited.Add(l)
			}
		}

		temp.Seek(0, 0)

		pageNum += maxWorkers

		tempMaxPage -= maxWorkers
		if tempMaxPage < maxWorkers {
			maxWorkers = tempMaxPage
		}
	}
}
