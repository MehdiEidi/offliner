package main

import (
	"bytes"
	"io"
	"net/http"
)

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
