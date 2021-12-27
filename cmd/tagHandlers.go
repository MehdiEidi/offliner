package main

import (
	"net/url"
	"path/filepath"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// handleImgTags is a handler called by the goquery Each method. It is for saving the images on the page. Also, it changes the tag's src so it references the local saved file.
func handleImgTags(index int, element *goquery.Selection, link string) {
	src, exists := element.Attr("src")
	if !exists {
		return
	}

	src, err := makeURLAbsolute(src, link)
	if err != nil {
		return
	}

	saveImg(src)

	src = makeURLRelative(src)
	element.SetAttr("src", src)
}

// handleScriptTags is a handler called by the goquery Each method. It is mainly for js files. It saves the js files on the disk. Also, it changes the tag's src so it references the local saved file.
func handleScriptTags(index int, element *goquery.Selection, link string) {
	src, exists := element.Attr("src")
	if !exists {
		return
	}

	if i := strings.LastIndex(src, "?"); i != -1 {
		src = src[:i]
	}

	src, err := makeURLAbsolute(src, link)
	if err != nil {
		return
	}

	ext := filepath.Ext(src)

	if ext == ".js" {
		saveScript(src)

		src = makeURLRelative(src)
		element.SetAttr("src", src)
	}
}

// handleLinkTags is a handler called by the goquery Each method. It is mainly for CSS files. It saves the css files on the disk. Also, it changes the tag's href so it references the local saved file.
func handleLinkTags(index int, element *goquery.Selection, link string) {
	href, exists := element.Attr("href")
	if !exists {
		return
	}

	if i := strings.LastIndex(href, "?"); i != -1 {
		href = href[:i]
	}

	href, err := makeURLAbsolute(href, link)
	if err != nil {
		return
	}

	ext := filepath.Ext(href)

	if ext == ".css" {
		saveCSS(href)

		href = makeURLRelative(href)
		element.SetAttr("href", href)
	}
}

// handleATags is a handler called by the goquery Each method. It checks some criteria about the found href attribute and adds it to the queue if it hasn't before.
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
	if href[len(href)-1] == '/' {
		href = href[:len(href)-1]
	}

	if strings.Contains(href, baseDomain) {
		if !visited.Has(href) {
			visited.Add(href)
			urls.Enqueue(href)

			href = makeURLRelative(href)
			element.SetAttr("href", href)
		}
	}
}
