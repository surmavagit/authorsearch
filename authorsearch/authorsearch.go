package authorsearch

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
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
	Data      []data
}

type searchResult struct {
	Resource string
	URL      string
	Error    string
}

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
			result.Error = err.Error()
		} else {
			result.URL = url
		}
		results = append(results, result)
	}
	return results, nil
}

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

func (website *Resource) loadCache() error {
	_, err := os.Stat(website.CacheFile)
	if errors.Is(err, os.ErrNotExist) {
		err := website.updateCache()
		if err != nil {
			return err
		}
	} else if err != nil {
		return err
	}

	file, err := os.ReadFile(website.CacheFile)
	if err != nil {
		return err
	}

	err = json.Unmarshal(file, &website.Data)
	return err
}

func (website Resource) updateCache() error {
	res, err := http.Get(website.BaseURL + website.QueryURL)
	if err != nil {
		return err
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}

	err = os.WriteFile(website.CacheFile, body, 0644)
	return err
}
