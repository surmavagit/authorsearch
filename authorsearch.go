package main

import (
	"strings"
)

type authorData struct {
	Description string `json:"name"`
	AuthorURL   string `json:"href"`
}

type resource struct {
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

type query struct {
	LastName  string
	FirstName string
	Year      string
}

// searchResource loads the cached data and searches for the author.
func (website resource) searchResource(query query, cacheDir string) resource {
	cacheFileName := cacheDir + "/" + website.Name + ".json"

	update, err := needUpdate(cacheFileName)
	if update {
		err = website.updateCache(cacheDir, cacheFileName)
	}
	if err != nil {
		website.Error = err
		return website
	}

	data, err := loadCache(cacheFileName)
	if err != nil {
		website.Error = err
		return website
	}

	results := []authorData{}
	for _, a := range data {
		if website.match(a.Description, query) {
			results = append(results, a)
		}
	}
	website.Results = results
	return website
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
