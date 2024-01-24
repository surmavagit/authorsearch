package authorsearch

import (
	"errors"
	"os"
	"strings"
)

type data struct {
	Description string `json:"name"`
	AuthorURL   string `json:"href"`
}

type Resource struct {
	Name      string
	BaseURL   string
	QueryURL  string
	CacheFile string
	URLFilter string
	Data      []data
}

type searchResult struct {
	Resource string
	URL      string
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
		url, err := r.SearchResource(query)
		if err != nil {
			result.ErrorMsg = err.Error()
		} else {
			result.URL = url
		}
		results = append(results, result)
	}
	return results, nil
}

// SearchResource loads the cached data and searches for the author.
// It returns the author URL on success and an empty string on failure.
func (website Resource) SearchResource(query string) (string, error) {
	err := website.loadCache()
	if err != nil {
		return "", err
	}

	for _, a := range website.Data {
		if strings.Contains(a.Description, query) {
			return website.BaseURL + a.AuthorURL, nil
		}
	}

	return "", nil
}
