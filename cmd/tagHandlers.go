package main

import (
	"bytes"
	"path/filepath"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

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

func handleImgTags(index int, element *goquery.Selection, link string) {
	src, exists := element.Attr("src")

	if exists {
		src = makeURLAbsolute(src, link)

		saveImg(src)

		src = makeURLRelative(src)
		element.SetAttr("src", src)
	}
}

func handleScriptTags(index int, element *goquery.Selection, link string) {
	src, exists := element.Attr("src")

	if exists {
		src = makeURLAbsolute(src, link)

		saveScript(src)

		src = makeURLRelative(src)
		element.SetAttr("src", src)
	}
}

func handleLinkTags(index int, element *goquery.Selection, link string) {
	href, exists := element.Attr("href")

	if exists {
		href = makeURLAbsolute(href, link)

		ext := filepath.Ext(href)

		if ext == ".css" {
			saveCSS(href)

			href = makeURLRelative(href)
			element.SetAttr("href", href)
		}
	}
}

func handleATags(index int, element *goquery.Selection, link string) {
	href, exists := element.Attr("href")

	if exists {
		href = makeURLAbsolute(href, link)

		if href[len(href)-1] == '/' {
			href = href[:len(href)-1]
		}

		if strings.Contains(href, baseDomain) {
			if !visited.Has(href) {
				visited.Add(href)
				urls.Enqueue(href)
			}
		}

		href = makeURLRelative(href)
		element.SetAttr("href", href)
	}
}
