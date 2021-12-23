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

	if link == "" || dirName == "" {
		log.Fatal("link or directory name cannot be empty")
	}

	temp, err := os.OpenFile("../temp/temp.txt", os.O_RDWR, 0644)
	if err != nil {
		log.Fatal(err)
	}

	processURL(link)

	// Join extracted URLs with a space into a string.
	lines := strings.Join(urls.Data, " ")

	temp.WriteString(lines + "\n")
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

	document, err := goquery.NewDocumentFromReader(bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	document.Find("a").Each(addLink)

	if err = save(bytes.NewBuffer(body), url); err != nil {
		return err
	}

	return nil
}

func addLink(index int, element *goquery.Selection) {
	href, exists := element.Attr("href")

	if href[len(href)-1] == '/' {
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

func save(body io.Reader, link string) error {
	filename, err := makeName(link)
	if err != nil {
		return err
	}

	file, err := os.Create("../output/" + dirName + "/pages/" + filename + ".html")
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
