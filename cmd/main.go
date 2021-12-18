package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"

	"github.com/MehdiEidi/offliner/workerpool"
	"github.com/PuerkitoBio/goquery"
)

var (
	// stack of URLs
	stack []string

	// set of visited URLs
	visited = make(map[string]bool)
	lock    sync.Mutex

	baseDomain string
)

func main() {
	homepage := flag.String("url", "", "URL of the homepage.")
	useProcess := flag.Bool("a", false, "If the concurrency must be done using multi processing instead of multi threading.")
	serial := flag.Bool("s", false, "If the crawling must be done in non-concurrent fashion.")
	maxPage := flag.Int("n", 100, "Max number of pages to save.")
	maxWorkers := flag.Int("p", 50, "Max number of concurrent execution units.")
	flag.Parse()

	if *homepage == "" {
		log.Fatal("homepage URL cannot be empty.")
	}

	findBase(*homepage)

	// create a separate directory for each link
	if _, err := os.Stat(baseDomain); os.IsNotExist(err) {
		err = os.Mkdir(baseDomain, os.ModePerm)
		if err != nil {
			log.Fatal(err)
		}
	}

	err := process(*homepage, -1)
	if err != nil {
		fmt.Println(err)
	}

	// serial,multi threading or multi processing ?
	switch {
	case *serial:
		fmt.Println("--serial processing--")

		for len(stack) != 0 {
			url := stack[len(stack)-1]
			stack = stack[:len(stack)-1]

			err := process(url, -1)
			if err != nil {
				fmt.Println(err)
			}
		}

	case *useProcess:

	default:
		fmt.Println("--multi threaded processing--")

		wp := workerpool.New(*maxWorkers, *maxPage, process)
		go wp.Run()

		for len(stack) != 0 {
			link := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			wp.AddTask(link)
		}
	}
}

func process(url string, workerID int) error {
	fmt.Println("worker", workerID, "Visiting:", url)

	lock.Lock()
	visited[url] = true
	lock.Unlock()

	res, err := http.Get(url)
	if err != nil {
		return err
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}

	res.Body = io.NopCloser(bytes.NewBuffer(body))

	document, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return err
	}
	document.Find("a").Each(linkHandler)

	res.Body = io.NopCloser(bytes.NewBuffer(body))

	err = save(res.Body, url)
	if err != nil {
		return err
	}

	return nil
}

func linkHandler(index int, element *goquery.Selection) {
	href, ok := element.Attr("href")

	if len(href) > 1 && href[len(href)-1] == '/' {
		href = href[:len(href)-1]
	}

	if ok && strings.Contains(href, baseDomain) {
		lock.Lock()
		_, ok := visited[href]
		lock.Unlock()
		if !ok {
			lock.Lock()
			stack = append(stack, href)
			lock.Unlock()
		}
	}
}

func save(body io.ReadCloser, url string) error {
	var filename string

	if strings.Contains(url, "http://") {
		filename = url[7:]
	} else {
		filename = url[8:]
	}

	fields := strings.Split(filename, "/")
	filename = strings.Join(fields, "-")

	file, err := os.Create(baseDomain + "/" + filename + ".html")
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, body)
	if err != nil {
		return err
	}

	return nil
}

func findBase(homepage string) {
	// remove trailing /
	if homepage[len(homepage)-1] == '/' {
		homepage = homepage[:len(homepage)-1]
	}

	// remove protocol scheme
	if strings.Contains(homepage, "http://") {
		baseDomain = homepage[7:]
	} else if strings.Contains(homepage, "https://") {
		baseDomain = homepage[8:]
	} else {
		baseDomain = homepage
	}
}
