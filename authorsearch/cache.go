package authorsearch

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"os"
	"time"
)

// loadCache creates the cache file for the resource, if one doesn't exist,
// and then loads the data into memory
func (website *Resource) loadCache() error {
	_, err := os.Stat(website.CacheFile)
	if errors.Is(err, os.ErrNotExist) {
		err = website.updateCache()
	}
	if err != nil {
		return err
	}

	file, err := os.ReadFile(website.CacheFile)
	if err != nil {
		return err
	}

	err = json.Unmarshal(file, &website.Data)
	return err
}

// updateCache carries out an http get request and saves the response body
// as a cache file
func (website Resource) updateCache() error {
	client := http.Client{
		Timeout: 5 * time.Second,
	}
	res, err := client.Get(website.BaseURL + website.QueryURL)
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

	err = os.WriteFile(website.CacheFile, body, 0644)
	return err
}
