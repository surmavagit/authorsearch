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
		output := r.Error + r.URL
		if output == "" {
			fmt.Println(r.Resource, "not found")
		} else {
			fmt.Println(r.Resource, output)
		}
	}
}

func checkInput(query []string) (string, error) {
	if len(query) <= 1 {
		return "", errors.New("no search query provided")
	}
	return strings.Join(query[1:], " "), nil
}
