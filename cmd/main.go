package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"

	"github.com/PuerkitoBio/goquery"
)

// var a = 1

var queue []string

var visited = make(map[string]bool)

// var visited sync.Map

var lock sync.RWMutex

var wg sync.WaitGroup

func main() {
	homepage := flag.String("url", "", "URL of the home page.")
	flag.Parse()

	if *homepage == "" {
		log.Fatal("home page URL cannot be empty")
	}

	// wg.Add(1)

	visited[*homepage] = true
	// visited.Store(*homepage, true)

	process(*homepage)

	for len(queue) != 0 {
		link := queue[0]
		queue = queue[1:]
		go process(link)
	}

	// time.Sleep(30 * time.Second)
	wg.Wait()
}

func process(link string) error {
	defer wg.Done()

	fmt.Println("Visiting:", link)
	lock.Lock()
	visited[link] = true
	lock.Unlock()
	// visited.Store(link, true)

	res, err := http.Get(link)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}

	res.Body = io.NopCloser(bytes.NewBuffer(body))

	document, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return err
	}
	document.Find("a").Each(getLink)

	res.Body = ioutil.NopCloser(bytes.NewBuffer(body))

	download(res.Body, link)

	fmt.Println("Finished:", link)
	return nil
}

func getLink(index int, element *goquery.Selection) {
	href, exists := element.Attr("href")
	if len(href) > 1 && href[len(href)-1] == '/' {
		href = href[:len(href)-1]
	}

	if exists {
		if strings.Contains(href, "urmia") {
			lock.RLock()
			defer lock.RUnlock()
			if _, ok := visited[href]; !ok {
				queue = append(queue, href)
				wg.Add(1)
			}
			// if _, ok := visited.Load(href); !ok {
			// 	queue = append(queue, href)
			// 	wg.Add(1)
			// }
		}
	}
}

// func content(url string) string {
// 	res, err := http.Get(url)
// 	if err != nil {
// 		log.Println("getContent - Request - ", err)
// 	}
// 	defer res.Body.Close()

// 	temp, _ := httputil.DumpResponse(res, true)
// 	b := bytes.NewBuffer(temp)

// 	document, _ := goquery.NewDocumentFromReader(res.Body)
// 	document.Find("a").Each(getLink)

// 	content, err := io.ReadAll(b)
// 	if err != nil {
// 		log.Println("getContent - ReadAll - ", err)
// 	}

// 	download(res.Body, url)

// 	return string(content)
// }

// func links(cont string) []string {
// 	re := regexp.MustCompile(`href="[[:graph:]^"]*"`)
// 	return re.FindAllString(cont, -1)[5:]
// }

func download(body io.ReadCloser, link string) error {
	URL, err := url.Parse(link)
	if err != nil {
		return err
	}
	fileName := URL.Query().Get("name")

	if fileName == "" {
		URLPath := URL.Path
		segments := strings.Split(URLPath, "/")
		fileName = segments[len(segments)-1]
	}

	if fileName == "" {
		// fileName = strconv.Itoa(a)
		// a++
		fileName = link[8 : len(link)-6]
	}

	file, err := os.Create("./pages/" + fileName + ".html")
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
