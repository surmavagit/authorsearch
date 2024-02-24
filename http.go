package main

import (
	"errors"
	"io"
	"net/http"
	"os"
	"time"
)

func requestURL(fullURL string) ([]byte, error) {
	client := http.Client{
		Timeout: 5 * time.Second,
	}
	res, err := client.Get(fullURL)

	if err != nil {
		return []byte{}, err
	}
	if res.StatusCode != 200 {
		return []byte{}, errors.New("Response status: " + res.Status)
	}

	body, err := io.ReadAll(res.Body)
	closeBody(res)
	return body, err
}

func closeBody(res *http.Response) {
	if res.Body == nil {
		return
	}

	err := res.Body.Close()
	if err != nil {
		os.Stderr.WriteString("Failed to close response body: " + err.Error())
		os.Exit(1)
	}
}
