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

	"github.com/PuerkitoBio/goquery"
)

var (
	// stack of URLs
	stack []string

	// set of visited URLs
	visited = make(map[string]bool)

	baseDomain string
)

func main() {
	homepage := flag.String("url", "", "URL of the homepage.")
	multiProcessing := flag.Bool("a", false, "If the concurrency must be done using multi processing instead of multi threading.")
	serial := flag.Bool("s", false, "If the crawling must be done in non-concurrent fashion.")
	// maxPage := flag.Int("n", math.MaxInt, "Max number of pages to crawl.")
	// maxConcurrency := flag.Int("p", 10, "Max number of concurrent execution units.")

	flag.Parse()

	if *homepage == "" {
		log.Fatal("homepage URL cannot be empty.")
	}

	baseDomain = findBase(*homepage)

	visited[*homepage] = true

	err := processHome(*homepage)
	if err != nil {
		fmt.Println(err)
	}

	// serial,multi threading or multi processing ?
	if *serial {
		fmt.Println("in serial")
		for len(stack) != 0 {
			url := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			err = serialProcess(url)
			if err != nil {
				fmt.Println(err)
			}
		}
	} else if *multiProcessing {

	} else {

	}
}

func processHome(homepage string) error {
	fmt.Println("Visiting:", homepage)

	res, err := http.Get(homepage)
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
	document.Find("a").Each(processLink)

	res.Body = io.NopCloser(bytes.NewBuffer(body))

	err = save(res.Body, homepage)
	if err != nil {
		return err
	}

	fmt.Println("Finished:", homepage)

	return nil
}

func processLink(index int, element *goquery.Selection) {
	href, ok := element.Attr("href")

	if len(href) > 1 && href[len(href)-1] == '/' {
		href = href[:len(href)-1]
	}

	if ok && strings.Contains(href, baseDomain) {
		if _, ok := visited[href]; !ok {
			stack = append(stack, href)
		}
	}
}

func findBase(homepage string) string {
	if homepage[len(homepage)-1] == '/' {
		homepage = homepage[:len(homepage)-1]
	}

	if strings.Contains(homepage, "http://") {
		return homepage[7:]
	} else if strings.Contains(homepage, "https://") {
		return homepage[8:]
	} else {
		return homepage
	}
}

func save(body io.ReadCloser, url string) error {
	var filename string

	if strings.Contains(url, "http://") {
		filename = url[7:]
	} else {
		filename = url[8:]
	}

	file, err := os.Create("./pages/" + filename + ".html")
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

func serialProcess(url string) error {
	fmt.Println("Visiting:", url)

	visited[url] = true

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
	document.Find("a").Each(processLink)

	res.Body = io.NopCloser(bytes.NewBuffer(body))

	err = save(res.Body, url)
	if err != nil {
		return err
	}

	fmt.Println("Finished:", url)

	return nil
}
