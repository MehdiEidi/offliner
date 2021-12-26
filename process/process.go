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
	urls       = stack.New()
	baseDomain string
)

func main() {
	link := os.Args[1]
	baseDomain = os.Args[2]
	connNum := os.Args[3]

	processURL(link)

	sendSocket, err := net.Dial("unix", "/tmp/ipc"+connNum+".sock")
	if err != nil {
		log.Fatal("In worker ", err)
	}
	defer sendSocket.Close()

	line := strings.Join(urls.Data, " ")

	sendSocket.Write([]byte(line))
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

	if exists && href != "" {
		if len(href) > 0 && href[len(href)-1] == '/' {
			href = href[:len(href)-1]
		}
		u, err := url.Parse(href)
		if err != nil {
			return
		}

		if u.Scheme == "" {
			return
		}

		if strings.Contains(href, baseDomain) {
			urls.Push(href)
		}
	}
}

func save(body io.Reader, link string) error {
	filename, err := makeName(link)
	if err != nil {
		return err
	}

	file, err := os.Create("../output/" + baseDomain + "/pages/" + filename + ".html")
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
