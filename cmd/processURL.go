package main

import (
	"bytes"
	"io"
	"net/http"

	"github.com/PuerkitoBio/goquery"
)

// processURL makes an https req to the given link, extracts its URLs, and saves the req body to the disk. The links are also edited on the page so they reference the local files.
func processURL(link string) error {
	res, err := http.Get(link)
	if err != nil {
		return err
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}

	body, err = extractLinks(body, link)
	if err != nil {
		return err
	}

	if err = savePage(bytes.NewBuffer(body), link); err != nil {
		return err
	}

	return nil
}

// extractLinks finds all the links on the page (body) and handles them via associated handlers. It makes a new html body out of the edited tags so the URLs on the page reference the local files. The static files are handled only if withfiles flag is true.
func extractLinks(body []byte, link string) ([]byte, error) {
	doc, err := goquery.NewDocumentFromReader(bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	doc.Find("a").Each(func(index int, element *goquery.Selection) {
		handleATags(index, element, link)
	})

	// Downloading static files happens only if withfiles flag is true.
	if flags.Withfiles {
		doc.Find("link").Each(func(index int, element *goquery.Selection) {
			handleLinkTags(index, element, link)
		})

		doc.Find("script").Each(func(index int, element *goquery.Selection) {
			handleScriptTags(index, element, link)
		})

		doc.Find("img").Each(func(index int, element *goquery.Selection) {
			handleImgTags(index, element, link)
		})
	}

	newBody, err := doc.Html()
	if err != nil {
		return nil, err
	}

	return []byte(newBody), nil
}
