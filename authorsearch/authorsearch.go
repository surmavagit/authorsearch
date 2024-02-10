package authorsearch

import (
	"errors"
	"fmt"
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
	Results    []authorData
	Error      error
}

var cacheFolder = "cache"

func init() {
	info, err := os.Stat(cacheFolder)

	if errors.Is(err, os.ErrNotExist) {
		err = os.Mkdir(cacheFolder, 0755)
		if err != nil {
			errMsg := fmt.Sprintf("can't create '%s' directory: %s", cacheFolder, err.Error())
			os.Stderr.WriteString(errMsg)
			os.Exit(1)
		}
		return
	}

	if err != nil {
		errMsg := fmt.Sprintf("can't access '%s' directory: %s", cacheFolder, err.Error())
		os.Stderr.WriteString(errMsg)
		os.Exit(1)
	}

	if !info.IsDir() {
		errMsg := fmt.Sprintf("'%s' is not a directory", cacheFolder)
		os.Stderr.WriteString(errMsg)
		os.Exit(1)
	}
}

// SearchResource loads the cached data and searches for the author.
func (website Resource) SearchResource(query string) Resource {
	data, err := website.loadCache()
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
