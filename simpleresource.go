package main

import (
	"strings"
)

type authorData struct {
	Description string `json:"name"`
	AuthorURL   string `json:"href"`
}

// searchResource loads the cached data and searches for the author.
func (website resource) searchResource(query query, cacheDir string) (data []authorData, err error) {
	if website.Complex {
		return website.searchComplexResource(query, cacheDir)
	}

	cacheFileName := cacheDir + "/" + website.Name + ".json"

	update, err := fileNotExist(cacheFileName)
	if err != nil {
		return []authorData{}, err
	}

	if update {
		data, err = website.getResource(cacheDir, cacheFileName)
	} else {
		err = loadFileJSON(cacheFileName, &data)
	}
	if err != nil {
		return []authorData{}, err
	}

	return website.filterRelevant(data, query), writeFileJSON(cacheFileName, data)
}

func (website resource) filterRelevant(data []authorData, query query) []authorData {
	results := []authorData{}
	for _, a := range data {
		if website.match(a.Description, query) {
			results = append(results, a)
		}
	}
	return results
}

func (website resource) match(authorDesc string, query query) bool {
	if !strings.Contains(authorDesc, query.LastName) {
		return false
	}

	if website.FirstName && !strings.Contains(authorDesc, query.FirstName) {
		return false
	}

	return !website.Year || strings.Contains(authorDesc, query.Year)
}

// updateCache carries out an http get request and saves the response body
// into a file
func (website resource) getResource(cacheDir string, cacheFileName string) ([]authorData, error) {
	fullURL := website.BaseURL + website.QueryURL
	body, err := requestURL(fullURL)
	if err != nil {
		return []authorData{}, err
	}

	data, err := website.readResource(body)
	if err != nil {
		return []authorData{}, err
	}

	return dedupe(data), nil
}
