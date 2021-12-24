package main

import (
	"bytes"
	"io"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/MehdiEidi/offliner/pkg/stack"
	"github.com/PuerkitoBio/goquery"
)

var (
	urls    *stack.Stack
	dirName string
)

func main() {
	link := os.Args[1]
	dirName = os.Args[2]
	connNum := os.Args[3]

	sendSocket, err := net.Dial("unix", "/tmp/ipc"+connNum+".sock")
	if err != nil {
		log.Fatal("In worker", err)
	}

	urls = stack.New()
	processURL(link)

	line := strings.Join(urls.Data, " ")

	sendSocket.Write([]byte(line))

	sendSocket.Close()
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
