package authorsearch

import (
	"errors"
	"os"
	"strings"
)

type authorData struct {
	Description string `json:"name"`
	AuthorURL   string `json:"href"`
}

type Resource struct {
	Name       string
	BaseURL    string
	QueryURL   string
	DataFormat string
	URLFilter  string // Valid URLs contain this string
}

type searchResult struct {
	Resource string
	Authors  []authorData
	ErrorMsg string
}

// Search creates a cache folder, if one doesn't exist, and then runs
// SearchResouce on every resource provided to it.
func Search(resource []Resource, query string) ([]searchResult, error) {
	_, err := os.Stat("cache")
	if errors.Is(err, os.ErrNotExist) {
		err = os.Mkdir("cache", 0755)
	}
	if err != nil {
		return []searchResult{}, err
	}

	results := []searchResult{}
	for _, r := range resource {
		result := searchResult{Resource: r.Name}
		dataSlice, err := r.SearchResource(query)
		if err != nil {
			result.ErrorMsg = err.Error()
		} else {
			result.Authors = dataSlice
		}
		results = append(results, result)
	}
	return results, nil
}

// SearchResource loads the cached data and searches for the author.
func (website Resource) SearchResource(query string) ([]authorData, error) {
	fullQueryURL := website.BaseURL + website.QueryURL
	cacheFileName := "cache/" + website.Name + "." + website.DataFormat
	cache, err := loadCache(fullQueryURL, cacheFileName)
	if err != nil {
		return []authorData{}, err
	}

	rawData, err := website.parseCache(cache)
	if err != nil {
		return []authorData{}, err
	}

	filteredData, err := website.filterData(rawData)
	if err != nil {
		return []authorData{}, err
	}

	results := []authorData{}
	for _, a := range filteredData {
		if strings.Contains(a.Description, query) {
			if strings.HasPrefix(a.AuthorURL, "/") {
				a.AuthorURL = website.BaseURL + a.AuthorURL
			} else {
				a.AuthorURL = website.BaseURL + "/" + a.AuthorURL
			}
			results = append(results, a)
		}
	}
	return results, nil
}
