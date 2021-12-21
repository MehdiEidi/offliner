package main

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/MehdiEidi/offliner/pkg/stack"
	"github.com/PuerkitoBio/goquery"
)

var (
	urls    = stack.New()
	dirName string
)

func main() {
	link := os.Args[1]
	dirName = os.Args[2]

	// var link string
	// scanner := bufio.NewScanner(os.Stdin)
	// scanner.Scan()
	// temp := scanner.Text()

	// t := strings.Split(temp, " ")

	// link = t[0]
	// dirName = t[1]

	temp, err := os.OpenFile("../temp/temp.txt", os.O_RDWR, 0644)
	if err != nil {
		log.Fatal(err)
	}

	processURL(link)

	allUrls := strings.Join(urls.Data, " ")

	temp.WriteString(allUrls + "\n")
}

func processURL(url string) error {
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
	document.Find("a").Each(addLink)

	res.Body = io.NopCloser(bytes.NewBuffer(body))

	err = save(res.Body, url)
	if err != nil {
		return err
	}

	return nil
}

func addLink(index int, element *goquery.Selection) {
	href, exists := element.Attr("href")

	// do this??
	if len(href) > 1 && href[len(href)-1] == '/' {
		href = href[:len(href)-1]
	}

	u, err := url.Parse(href)
	if err == nil {
		if u.Scheme == "" {
			return
		}
	}

	if exists && strings.Contains(href, dirName) {
		urls.Push(href)
	}
}

// save saves the body of the given URL to a file.
func save(body io.ReadCloser, link string) error {
	filename, err := makeName(link)
	if err != nil {
		return err
	}

	file, err := os.Create(dirName + "/" + filename + ".html")
	if err != nil {
		return err
	}
	defer file.Close()

	if _, err = io.Copy(file, body); err != nil {
		return err
	}

	return nil
}

func makeName(link string) (filename string, err error) {
	u, err := url.Parse(link)
	if err != nil {
		return "", err
	}
	u.Scheme = ""

	filename = u.String()[2:]
	filename = strings.ReplaceAll(filename, "/", "_")

	return filename, nil
}
