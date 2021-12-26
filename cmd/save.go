package main

import (
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
)

func savePage(body io.Reader, link string) error {
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

	// Updating the number of the saved pages and progress bar.
	progress.Add(1)

	return nil
}

func saveCSS(link string) error {
	res, err := http.Get(link)
	if err != nil {
		return err
	}

	filename, err := makeName(link)
	if err != nil {
		return err
	}

	file, err := os.Create("../output/" + baseDomain + "/static/css/" + filename)
	if err != nil {
		return err
	}
	defer file.Close()

	if _, err = io.Copy(file, res.Body); err != nil {
		return err
	}

	return nil
}

func saveScript(link string) error {
	res, err := http.Get(link)
	if err != nil {
		return err
	}

	filename, err := makeName(link)
	if err != nil {
		return err
	}

	file, err := os.Create("../output/" + baseDomain + "/static/js/" + filename)
	if err != nil {
		return err
	}
	defer file.Close()

	if _, err = io.Copy(file, res.Body); err != nil {
		return err
	}

	return nil
}

func saveImg(link string) error {
	res, err := http.Get(link)
	if err != nil {
		return err
	}

	filename, err := makeName(link)
	if err != nil {
		return err
	}

	file, err := os.Create("../output/" + baseDomain + "/static/img/" + filename)
	if err != nil {
		return err
	}
	defer file.Close()

	if _, err = io.Copy(file, res.Body); err != nil {
		return err
	}

	return nil
}

func saveFile(link string) error {
	res, err := http.Get(link)
	if err != nil {
		return err
	}

	filename, err := makeName(link)
	if err != nil {
		return err
	}

	file, err := os.Create("../output/" + baseDomain + "/static/files/" + filename)
	if err != nil {
		return err
	}
	defer file.Close()

	if _, err = io.Copy(file, res.Body); err != nil {
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

	filename = u.String()[2:]
	filename = strings.ReplaceAll(filename, "/", "_")

	return filename, nil
}
