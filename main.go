package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

type data struct {
	Description string `json:"name"`
	AuthorURL   string `json:"href"`
}

type resource struct {
	baseURL   string
	queryURL  string
	cacheFile string
	data      []data
}

var marxists = resource{
	baseURL:   "https://www.marxists.org/",
	queryURL:  "admin/js/data/authors.json",
	cacheFile: "cache/marxists.json",
}

func init() {
	_, err := os.Stat("cache")
	if errors.Is(err, os.ErrNotExist) {
		err = os.Mkdir("cache", 0755)
	}
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func main() {
	if len(os.Args) == 1 {
		fmt.Println("no search query provided")
		os.Exit(1)
	}
	searchQuery := strings.Join(os.Args[1:], " ")

	err := marxists.loadCache()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	authorLink := marxists.search(searchQuery)
	if authorLink == "" {
		fmt.Println("not found")
		os.Exit(1)
	}
	fmt.Println(authorLink)
}

func (website resource) search(query string) string {
	for _, a := range website.data {
		if strings.Contains(a.Description, query) {
			return website.baseURL + a.AuthorURL
		}
	}

	return ""
}

func (website *resource) loadCache() error {
	_, err := os.Stat(website.cacheFile)
	if errors.Is(err, os.ErrNotExist) {
		err := website.updateCache()
		if err != nil {
			return err
		}
	} else if err != nil {
		return err
	}

	file, err := os.ReadFile(website.cacheFile)
	if err != nil {
		return err
	}

	err = json.Unmarshal(file, &website.data)
	return err
}

func (website resource) updateCache() error {
	res, err := http.Get(website.baseURL + website.queryURL)
	if err != nil {
		return err
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}

	err = os.WriteFile(website.cacheFile, body, 0644)
	return err
}
