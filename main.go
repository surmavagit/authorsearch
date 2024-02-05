package main

import (
	"errors"
	"fmt"
	"os"
	"strings"
)

func main() {
	searchQuery, err := checkInput(os.Args)
	if err != nil {
		os.Stderr.WriteString(err.Error())
		os.Exit(1)
	}

	for _, r := range resources {
		dataSlice, err := r.SearchResource(searchQuery)

		if err != nil {
			fmt.Printf("%-10s  %s\n", r.Name, err.Error())
			continue
		}

		numLinks := len(dataSlice)
		if numLinks == 0 {
			fmt.Printf("%-10s  [0 of 0]  not found\n", r.Name)
			continue
		}

		for i, l := range dataSlice {
			fmt.Printf("%-10s  [%d of %d]  %-30s  %s\n", r.Name, i+1, numLinks, l.Description, l.AuthorURL)
		}
	}
}

func checkInput(query []string) (string, error) {
	if len(query) <= 1 {
		return "", errors.New("no search query provided")
	}
	return strings.Join(query[1:], " "), nil
}
