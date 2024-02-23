package main

import (
	"encoding/json"
	"errors"
	"os"
)

// fileNotExist returns true and error if the cache file with the specified
// name doesn't exist. It returns false if either there is no error or
// if there is a different error
func fileNotExist(fileName string) (bool, error) {
	_, err := os.Stat(fileName)
	return errors.Is(err, os.ErrNotExist), err
}

// loadFileJSON reads the cache file and unmarshals the json into the target
func loadFileJSON(fileName string, target interface{}) error {
	stream, err := os.ReadFile(fileName)
	if err != nil {
		return err
	}

	err = json.Unmarshal(stream, target)
	return err
}

// writeFileJSON turns source into indented JSON and writes it into file
func writeFileJSON(fileName string, source any) error {
	stream, err := json.MarshalIndent(source, "", "    ")
	if err != nil {
		return err
	}

	return os.WriteFile(fileName, []byte(stream), 0644)
}

// updateCache carries out an http get request and saves the response body
// into a file
func (website resource) updateCache(cacheDir string, cacheFileName string) ([]authorData, error) {
	fullURL := website.BaseURL + website.QueryURL
	body, err := getResource(fullURL)
	if err != nil {
		return []authorData{}, err
	}

	data, err := website.readResource(body)
	if err != nil {
		return []authorData{}, err
	}
	filteredData := website.dedupe(data)

	return filteredData, writeFileJSON(cacheFileName, filteredData)
}
