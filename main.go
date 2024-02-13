package main

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/surmavagit/authorsearchcli/authorsearch"
)

func main() {
	searchQuery, err := checkInput(os.Args)
	if err != nil {
		os.Stderr.WriteString(err.Error())
		os.Exit(1)
	}

	cacheDirectory := "cache"
	err = checkCacheDir(cacheDirectory)
	if err != nil {
		os.Stderr.WriteString(err.Error())
		os.Exit(1)
	}

	dataChan := make(chan authorsearch.Resource)
	wg := sync.WaitGroup{}
	for _, r := range resources {
		resource := r
		wg.Add(1)
		go func() {
			dataChan <- resource.SearchResource(searchQuery, cacheDirectory)
		}()
	}
	go func() {
		for r := range dataChan {
			printResults(r)
			wg.Done()
		}
	}()
	wg.Wait()
}

func printResults(r authorsearch.Resource) {
	if r.Error != nil {
		fmt.Printf("%-10s  %s\n", r.Name, r.Error.Error())
		return
	}

	numLinks := len(r.Results)
	if numLinks == 0 {
		fmt.Printf("%-10s  [0 of 0]  not found\n", r.Name)
		return
	}

	for i, l := range r.Results {
		fmt.Printf("%-10s  [%d of %d]  %-35s  %s\n", r.Name, i+1, numLinks, l.Description, l.AuthorURL)
	}
}

func checkInput(query []string) (string, error) {
	if len(query) <= 1 {
		return "", errors.New("no search query provided")
	}
	return strings.Join(query[1:], " "), nil
}

func checkCacheDir(directory string) error {
	info, err := os.Stat(directory)

	if errors.Is(err, os.ErrNotExist) {
		err = os.Mkdir(directory, 0755)
		if err != nil {
			return err
		}
		return nil
	}

	if err != nil {
		return err
	}

	if !info.IsDir() {
		return fmt.Errorf("%s is not a directory", directory)
	}

	return nil
}
