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

var cacheFolder = "cache"

func init() {
	info, err := os.Stat(cacheFolder)
	if errors.Is(err, os.ErrNotExist) {
		err = os.Mkdir(cacheFolder, 0755)
	}
	if err != nil || !info.IsDir() {
		os.Stderr.WriteString("can't access cache folder")
		os.Exit(1)
	}
}

// SearchResource loads the cached data and searches for the author.
func (website Resource) SearchResource(query string) ([]authorData, error) {
	fullQueryURL := website.BaseURL + website.QueryURL
	cacheFileName := cacheFolder + "/" + website.Name + "." + website.DataFormat
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
