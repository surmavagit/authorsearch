package authorsearch

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

// loadCache reads the cache file, filters the data, if necessary,
// and loads the results into memory. If the cache file doesn't exist,
// loadCache runs updateCache function.
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

	var rawData []data
	err = json.Unmarshal(file, &rawData)
	if err != nil {
		return err
	}

	if website.URLFilter != "" {
		for _, d := range rawData {
			if strings.Contains(d.AuthorURL, website.URLFilter) {
				website.Data = append(website.Data, d)
			}
		}
	} else {
		website.Data = rawData
	}
	return nil
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
