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

	file, err := os.ReadFile(cacheFileName)
	if err != nil {
		return []byte{}, err
	}
	return file, nil
}

// updateCache carries out an http get request and saves the response body
// into a file
func updateCache(fullURL string, cacheFileName string) error {
	client := http.Client{
		Timeout: 5 * time.Second,
	}
	res, err := client.Get(fullURL)
	if err != nil {
		return err
	}
	if res.StatusCode != 200 {
		return errors.New(res.Status)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}

	err = os.WriteFile(cacheFileName, body, 0644)
	return err
}
