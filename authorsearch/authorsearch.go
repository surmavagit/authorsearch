package authorsearch

import (
	"strings"
)

type authorData struct {
	Description string `json:"name"`
	AuthorURL   string `json:"href"`
}

type Resource struct {
	Name         string
	BaseURL      string
	QueryURL     string
	DataFormat   string
	DescInParent bool
	URLFilter    string // Valid URLs contain this string
	Results      []authorData
	Error        error
}

// SearchResource loads the cached data and searches for the author.
func (website Resource) SearchResource(query string, cacheDir string) Resource {
	data, err := website.loadCache(cacheDir)
	if err != nil {
		website.Error = err
		return website
	}

	results := []authorData{}
	for _, a := range data {
		if search(a.Description, query) {
			results = append(results, a)
		}
	}
	website.Results = results
	return website
}

func search(authorDesc string, query string) bool {
	return strings.Contains(authorDesc, query)
}
