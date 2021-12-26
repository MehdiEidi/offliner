package main

import (
	"net/url"
	"path/filepath"
)

func makeURLAbsolute(href string, currentPageURL string) string {
	// Empty href references the current page.
	if href == "" {
		return currentPageURL
	}

	u, err := url.Parse(href)
	if err == nil {
		if u.Scheme == "" {
			if href[0] == '/' { // Relative to root.
				if rootURL[len(rootURL)-1] == '/' {
					href = rootURL + href[1:]
				} else {
					href = rootURL + href
				}
			} else { // Relative to current page.
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
	}

	return href
}

func makeURLRelative(href string) string {
	filename, _ := makeName(href)
	filename = url.QueryEscape(filename)

	ext := filepath.Ext(filename)

	switch ext {
	case ".css":
		href = "../static/css/" + filename

	case ".js":
		href = "../static/js/" + filename

	case ".jpg", ".png":
		href = "../static/img/" + filename

	case ".html", ".htm":
		href = "./" + filename

	default:
		href = "../static/files/" + filename
	}

	return href
}
