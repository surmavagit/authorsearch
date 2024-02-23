package main

import (
	"encoding/json"
	"errors"
	"os"
)

// needUpdate returns true and error if the cache file with the specified
// name doesn't exist. It returns false if either there is no error or
// if there is a different error
func needUpdate(cacheFileName string) (bool, error) {
	_, err := os.Stat(cacheFileName)
	return errors.Is(err, os.ErrNotExist), err
}

// loadCache reads the cache file and unmarshals the json
func loadCache(cacheFileName string) ([]authorData, error) {
	stream, err := os.ReadFile(cacheFileName)
	if err != nil {
		return []authorData{}, err
	}

	data := []authorData{}
	err = json.Unmarshal(stream, &data)
	return data, err
}

// updateCache carries out an http get request and saves the response body
// into a file
func (website resource) updateCache(cacheDir string, cacheFileName string) error {
	fullURL := website.BaseURL + website.QueryURL
	body, err := getResource(fullURL)
	if err != nil {
		return err
	}

	data, err := website.readResource(body)
	if err != nil {
		return err
	}
	filteredData := website.filterAndDedupe(data)

	stream, err := json.MarshalIndent(filteredData, "", "    ")
	if err != nil {
		return err
	}

	err = os.WriteFile(cacheFileName, []byte(stream), 0644)
	return err
}
