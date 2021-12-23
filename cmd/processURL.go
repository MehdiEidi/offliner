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

// processURL gets a URL. It makes an http get request to the URL. It basically does two things with the response body, first it extracts all the links on the page and adds them to the stack, second, it saves the page on the disk.
func processURL(link string, workerID int) error {
	res, err := http.Get(link)
	if err != nil {
		return err
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}

	// Extracting links.
	doc, err := goquery.NewDocumentFromReader(bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	doc.Find("a").Each(handleLink)

	if err = savePage(bytes.NewBuffer(body), link); err != nil {
		return err
	}

	return nil
}

// handleLink checks some criteria about the found link, makes relative paths absolute, and adds them to the stack and set.
func handleLink(index int, element *goquery.Selection) {
	href, exists := element.Attr("href")

	if exists && href != "" {
		// Make relative paths absolute.
		u, err := url.Parse(href)
		if err == nil {
			if u.Scheme == "" {
				if href[0] == '/' {
					if homeURL[len(homeURL)-1] == '/' {
						href = homeURL + href[1:]
					} else {
						href = homeURL + href
					}
				} else {
					if homeURL[len(homeURL)-1] == '/' {
						href = homeURL + href
					} else {
						href = homeURL + "/" + href
					}
				}
			}
		}

		if strings.Contains(href, baseDomain) {
			if !visited.Has(href) {
				visited.Add(href)
				urls.Push(href)
			}
		}
	}
}

// savePage saves the body of the given URL to a file.
func savePage(body io.Reader, link string) error {
	filename, err := makeName(link)
	if err != nil {
		return err
	}

	file, err := os.Create(baseDomain + "/pages/" + filename + ".html")
	if err != nil {
		return err
	}
	defer file.Close()

	if _, err = io.Copy(file, body); err != nil {
		return err
	}

	// Updating the number of the saved pages and progress bar.
	progress.Add(1)

	return nil
}

// makeName constructs a savable name for a file out of the given URL.
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

// createDirs contains a slice of the directories to be created, then it ranges over the slice and creates those directories.
func createDirs() error {
	// Slice of the directories to be created. Order is important.
	dirs := []string{baseDomain, baseDomain + "/pages/", baseDomain + "/static/", baseDomain + "/static/css/", baseDomain + "/static/js/", baseDomain + "/static/files/"}

	for _, d := range dirs {
		if _, err := os.Stat(d); os.IsNotExist(err) {
			if err = os.Mkdir(d, os.ModePerm); err != nil {
				return err
			}
		}
	}

	return nil
}

// setHomeURL sets the homeURL global variable to the full URL of the homepage.
func setHomeURL(link string) error {
	u, err := url.Parse(link)
	if err != nil {
		return err
	}
	homeURL = u.Scheme + u.Host
	return nil
}

// setBaseDomain sets the baseDomain to the Host field of the homepage URL. It returns errors if any.
func setBaseDomain(homepage string) error {
	u, err := url.Parse(homepage)
	if err != nil {
		return err
	}
	baseDomain = u.Host
	return nil
}
