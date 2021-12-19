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

func process(url string, workerID int) error {
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
	href, exists := element.Attr("href")

	// do this??
	if len(href) > 1 && href[len(href)-1] == '/' {
		href = href[:len(href)-1]
	}

	u, _ := url.Parse(href)
	if u.Scheme == "" {
		return
	}

	if exists && strings.Contains(href, baseDomain) {
		if !visited.Has(href) {
			visited.Add(href)
			urls.Push(href)
		}
	}
}

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
