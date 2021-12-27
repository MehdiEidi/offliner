package main

import (
	"net/url"
	"path/filepath"
	"strings"
)

// makeURLAbsolute gets a link which might be a relative URL. It converts the relative to an absolute URL.
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
			// if currentPageURL[len(currentPageURL)-1] == '/' {
			// 	href = currentPageURL + href
			// } else {
			// 	href = currentPageURL + "/" + href
			// }
			if rootURL[len(rootURL)-1] == '/' {
				href = rootURL + href
			} else {
				href = rootURL + "/" + href
			}
		}
	}

	return href, nil
}

// makeURLRelative gets an absolute URL and converts it to a relative URL which will be used to reference the file locally.
func makeURLRelative(href string) string {
	filename, _ := makeName(href)
	filename = url.QueryEscape(filename)

	ext := filepath.Ext(filename)

	switch ext {
	case ".css":
		href = "../static/css/" + filename

	case ".js":
		href = "../static/js/" + filename

	case ".jpg", ".png", ".jpeg":
		href = "../static/img/" + filename

	default:
		href = "./" + filename
	}

	return href
}

// removeAnchor gets a link and removes its anchor.
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
