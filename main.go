package main

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/surmavagit/authorsearchcli/authorsearch"
)

func main() {
	searchQuery, err := checkInput(os.Args)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	results, err := authorsearch.Search(resources, searchQuery)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	for _, r := range results {

		if r.ErrorMsg != "" {
			fmt.Println(r.Resource, r.ErrorMsg)
			continue
		}

		numLinks := len(r.Authors)
		if numLinks == 0 {
			fmt.Printf("%-10s  [0 of 0]  not found\n", r.Resource)
			continue
		}
		for i, l := range r.Authors {
			fmt.Printf("%-10s  [%d of %d]  %-30s  %s\n", r.Resource, i+1, numLinks, l.Description, l.AuthorURL)
		}
	}
}

func checkInput(query []string) (string, error) {
	if len(query) <= 1 {
		return "", errors.New("no search query provided")
	}
	return strings.Join(query[1:], " "), nil
}
