package main

import (
	"bytes"
	"io"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
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

	sendSocket, err := net.Dial("unix", "/tmp/ipc"+connNum+".sock")
	if err != nil {
		log.Fatal("In worker ", err)
	}
	defer sendSocket.Close()

	processURL(link)

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

	document.Find("a").Each(func(index int, element *goquery.Selection) {
		handleATags(index, element, url)
	})

	if err = save(bytes.NewBuffer(body), url); err != nil {
		return err
	}

	return nil
}

func handleATags(index int, element *goquery.Selection, link string) {
	href, exists := element.Attr("href")
	if !exists {
		return
	}

	if strings.HasPrefix(href, "data:") || strings.HasPrefix(href, "mailto:") {
		return
	}

	href = removeAnchor(href)

	u, err := url.Parse(href)
	if err != nil {
		return
	}

	if u.Port() != "" {
		return
	}

	if !isHTML(href) {
		return
	}

	// Not handling query params.
	if strings.Contains(href, "?") {
		return
	}

	href, err = makeURLAbsolute(href, link)
	if err != nil {
		return
	}

	// Removing trailing slash.
	if len(href) > 0 && href[len(href)-1] == '/' {
		href = href[:len(href)-1]
	}

	if strings.Contains(href, baseDomain) {
		urls.Push(href)
	}
}

func save(body io.Reader, link string) error {
	filename, err := makeName(link)
	if err != nil {
		return err
	}

	file, err := os.Create("../output/" + baseDomain + "/pages/" + filename)
	if err != nil {
		return err
	}
	defer file.Close()

	if _, err = io.Copy(file, body); err != nil {
		return err
	}

	return nil
}

// makeName constructs a savable name for a file out of the given URL.
func makeName(link string) (filename string, err error) {
	u, err := url.Parse(link)
	if err != nil {
		return "", err
	}

	u.Scheme = ""

	filename = u.String()[2:] // Skip slashes.
	filename = strings.ReplaceAll(filename, "/", "_")

	return filename, nil
}

func makeURLAbsolute(href string, currentPageURL string) (string, error) {
	rootURL, err := getRoot(currentPageURL)
	if err != nil {
		return "", nil
	}

	// Empty href references the current page.
	if href == "" {
		return currentPageURL, nil
	}

	u, err := url.Parse(href)
	if err != nil {
		return "", err
	}

	if u.Scheme == "" {
		if href[0] == '/' {
			if rootURL[len(rootURL)-1] == '/' {
				href = rootURL + href[1:]
			} else {
				href = rootURL + href
			}
		} else {
			if rootURL[len(rootURL)-1] == '/' {
				href = rootURL + href
			} else {
				href = rootURL + "/" + href
			}
		}
	}

	return href, nil
}

func removeAnchor(href string) string {
	i := strings.LastIndex(href, "#")
	if i == -1 {
		return href
	}
	return href[:i]
}

// isHTML checks to see if the given link is referencing an HTML page.
func isHTML(href string) bool {
	i := strings.LastIndex(href, "/")

	ext := filepath.Ext(href[i+1:])

	if ext != "" && ext != ".html" {
		return false
	}

	return true
}

// getRoot returns the root (host) of the given URL.
func getRoot(link string) (string, error) {
	u, err := url.Parse(link)
	if err != nil {
		return "", nil
	}
	return u.Scheme + "://" + u.Host, nil
}
