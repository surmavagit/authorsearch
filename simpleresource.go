package main

import (
	"strings"
)

type authorData struct {
	Description string `json:"name"`
	AuthorURL   string `json:"href"`
}

// searchResource loads the cached data and searches for the author.
func (website resource) searchResource(q query, cacheDir string) (data []authorData, cacheErr error, err error) {
	cacheFileName := website.getCacheFileName(cacheDir, q.LastName)
	data, history, cacheErr := website.searchInCache(q, cacheFileName)
	if data != nil {
		// successfully found data in cache
		return data, nil, nil
	}

	data, err = website.getResource(q)
	if err != nil {
		// didn't find data in cache, error getting data from resource
		return nil, cacheErr, err
	}
	filteredData := website.filterRelevant(data, q)

	// don't update cache if there already is a cache error
	if cacheErr == nil {
		if website.Complex {
			if history == nil {
				history = records{}
			}
			history[getQueryString(q)] = filteredData
			cacheErr = writeFileJSON(cacheFileName, history)
		} else {
			cacheErr = writeFileJSON(cacheFileName, data)
		}
	}

	return filteredData, cacheErr, nil
}

func (website resource) searchInCache(q query, cacheFileName string) (data []authorData, history records, err error) {
	noFile, err := fileNotExist(cacheFileName)
	if err != nil || noFile {
		return nil, nil, err
	}

	if !website.Complex {
		err = loadFileJSON(cacheFileName, &data)
	} else {
		err = loadFileJSON(cacheFileName, &history)
	}
	if err != nil {
		return nil, nil, err
	}

	if website.Complex {
		found, data := searchInHistory(history, q)
		if !found {
			return nil, history, nil
		}
		return data, nil, nil
	}
	return website.filterRelevant(data, q), nil, nil
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
	if website.QueryFirst && query.FirstName != "" {
		queryTerms = append(queryTerms, query.FirstName)
	}
	if website.QueryYear && query.Year != "" {
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
