package main

import (
	"bytes"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// processURL gets a URL and a worker ID. It makes a http get request to the URL. It basically does two things with the response body, first it extracts all the links on the page and adds them to the stack, second, it saves the page on the disk.
func processURL(url string, workerID int) error {
	res, err := http.Get(url)
	if err != nil {
		return err
	}

	// Because we need to read the body twice (once for extracting links and once for saving the page), we first need to save the body in a []byte. Later we can use this slice to rewind the body.
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}

	// NewBuffer returns a Reader value and NopCloser just adds a Closer method on it, making this value a ReadCloser which can get assigned back to the res.Body. In that case we can read the body again.
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

// addLink checks some criteria about the found link and adds it to the stack and set.
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

	if exists && strings.Contains(href, baseDomain) {
		if !visited.Has(href) {
			visited.Add(href)
			urls.Push(href)
		}
	}
}

// save saves the body of the given URL to a file.
func save(body io.ReadCloser, link string) error {
	filename, err := makeName(link)
	if err != nil {
		return err
	}

	file, err := os.Create(baseDomain + "/" + filename + ".html")
	if err != nil {
		return err
	}
	defer file.Close()

	if _, err = io.Copy(file, body); err != nil {
		return err
	}

	// update the number of the saved pages and progress bar.
	pageNum++
	bar.Add(1)

	return nil
}

// findBaseDomain sets the baseDomain to the Host field of the homepage URL. It returns errors if any.
func findBaseDomain(homepage string) error {
	u, err := url.Parse(homepage)
	if err != nil {
		return err
	}

	baseDomain = u.Host

	return nil
}

// createDir creates a directory for the baseDomain if not exists.
func createDir() error {
	if _, err := os.Stat(baseDomain); os.IsNotExist(err) {
		err = os.Mkdir(baseDomain, os.ModePerm)
		if err != nil {
			return err
		}
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
