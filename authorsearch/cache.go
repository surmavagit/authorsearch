package authorsearch

import (
	"errors"
	"io"
	"net/http"
	"os"
	"time"
)

// loadCache loads the contents of the cache file. If it doesn't
// exist, updateCache function is called.
func loadCache(fullURL string, cacheFileName string) ([]byte, error) {
	_, err := os.Stat(cacheFileName)
	if errors.Is(err, os.ErrNotExist) {
		err = updateCache(fullURL, cacheFileName)
	}
	if err != nil {
		return []byte{}, err
	}

	return os.ReadFile(cacheFileName)
}

// updateCache carries out an http get request and saves the response body
// into a file
func updateCache(fullURL string, cacheFileName string) error {
	body, err := getResource(fullURL)
	if err != nil {
		return err
	}
	err = os.WriteFile(cacheFileName, body, 0644)
	return err
}

func getResource(fullURL string) ([]byte, error) {
	client := http.Client{
		Timeout: 5 * time.Second,
	}
	res, err := client.Get(fullURL)
	defer closeConnection(res)

	if err != nil {
		return []byte{}, err
	}
	if res.StatusCode != 200 {
		return []byte{}, errors.New("Response status: " + res.Status)
	}

	body, err := io.ReadAll(res.Body)
	return body, err
}

func closeConnection(res *http.Response) {
	err := res.Body.Close()
	if err != nil {
		os.Stderr.WriteString("Failed to close connection: " + err.Error())
		os.Exit(1)
	}
}
