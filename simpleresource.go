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
	cacheFileName := website.getCacheFileName(cacheDir, query.LastName)
	if website.Complex {
		return website.searchComplexResource(query, cacheFileName)
	}

	update, err := fileNotExist(cacheFileName)
	if err != nil {
		return []authorData{}, err
	}

	if update {
		data, err = website.getResource(query)
	} else {
		err = loadFileJSON(cacheFileName, &data)
	}
	if err != nil {
		return []authorData{}, err
	}

	return website.filterRelevant(data, query), writeFileJSON(cacheFileName, data)
}

func (website resource) getCacheFileName(cacheDir string, lastName string) string {
	if !website.Complex {
		return cacheDir + "/" + website.Name + ".json"
	}
	return cacheDir + "/" + website.Name + "_" + strings.ToLower(lastName) + ".json"
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

	if website.Year && !strings.Contains(authorDesc, query.Year) {
		return false
	}

	return true
}

func (website resource) getFullQueryURL(query query) string {
	if !website.Complex {
		return website.BaseURL + website.QueryURL
	}

	queryTerms := []string{query.LastName}
	if website.FirstName && query.FirstName != "" {
		queryTerms = append(queryTerms, query.FirstName)
	}
	if website.Year && query.Year != "" {
		queryTerms = append(queryTerms, query.Year)
	}
	formatQuery := strings.Join(queryTerms, "+")

	return website.BaseURL + website.QueryURL + formatQuery
}

// getResource carries out an http get request and gets all the data
// from the response body
func (website resource) getResource(query query) ([]authorData, error) {
	body, err := requestURL(website.getFullQueryURL(query))
	if err != nil {
		return []authorData{}, err
	}
	data, err := website.readResource(body)
	if err != nil {
		return []authorData{}, err
	}

	// data should be deduped only in case of simple resources, where
	// creators often provide lists with duplicate values
	if website.Complex {
		return data, nil
	}

	return dedupe(data), nil
}
