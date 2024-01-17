package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
)

type data struct {
	Description string `json:"name"`
	AuthorURL   string `json:"href"`
}

type resource struct {
	baseURL      string
	queryURL     string
	sourceType   string
	resultFormat string
	data         []data
}

var marxists = resource{
	baseURL:  "https://www.marxists.org/",
	queryURL: "admin/js/data/authors.json",
}

func main() {
	if len(os.Args) == 1 {
		fmt.Println("no search query provided")
		os.Exit(1)
	}
	searchQuery := strings.Join(os.Args[1:], " ")

	err := marxists.loadDB()
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

func (website *resource) loadDB() error {
	res, err := http.Get(website.baseURL + website.queryURL)
	if err != nil {
		return err
	}

	err = json.NewDecoder(res.Body).Decode(&website.data)
	if err != nil {
		return err
	}

	return nil
}
