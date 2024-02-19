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
	DataFormat   string // json or html
	DescInParent bool   // Location of the full description in HTML document relative to the <a> tag
	URLFilter    string // Valid URLs contain this string
	FirstName    bool
	Year         bool
	Results      []authorData
	Error        error
}

type Query struct {
	LastName  string
	FirstName string
	Year      string
}

// SearchResource loads the cached data and searches for the author.
func (website Resource) SearchResource(query Query, cacheDir string) Resource {
	data, err := website.loadCache(cacheDir)
	if err != nil {
		website.Error = err
		return website
	}

	results := []authorData{}
	for _, a := range data {
		if website.search(a.Description, query) {
			results = append(results, a)
		}
	}
	website.Results = results
	return website
}

func (website Resource) search(authorDesc string, query Query) bool {
	if !strings.Contains(authorDesc, query.LastName) {
		return false
	}

	if website.FirstName && !strings.Contains(authorDesc, query.FirstName) {
		return false
	}

	return !website.Year || strings.Contains(authorDesc, query.Year)
}
