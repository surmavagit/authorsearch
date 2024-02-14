package main

import (
	"errors"
	"fmt"
	"os"
	"regexp"
	"sync"

	"github.com/surmavagit/authorsearchcli/authorsearch"
)

func main() {
	searchQuery, err := checkInput(os.Args[1:])
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
		os.Stderr.WriteString(fmt.Sprintf("%-10s  %s\n", r.Name, r.Error.Error()))
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

func checkInput(query []string) (authorsearch.Query, error) {
	if len(query) == 0 {
		return authorsearch.Query{}, errors.New("no search query provided")
	}

	if len(query) > 3 {
		return authorsearch.Query{}, errors.New("too many arguments")
	}

	queryStruct := authorsearch.Query{}
	for _, a := range query {
		numeric, err := regexp.MatchString("\\d", a)
		if err != nil {
			return authorsearch.Query{}, err
		}

		nonNumeric, err := regexp.MatchString("\\D", a)
		if err != nil {
			return authorsearch.Query{}, err
		}

		if numeric && nonNumeric {
			return authorsearch.Query{}, errors.New("invalid argument: numeric and nonnumeric characters")
		}

		if numeric && queryStruct.Year != "" {
			return authorsearch.Query{}, errors.New("only one year can be specified")
		}

		if numeric {
			queryStruct.Year = a
			continue
		}

		if queryStruct.LastName == "" {
			queryStruct.LastName = a
			continue
		}

		if queryStruct.FirstName == "" {
			queryStruct.FirstName = a
			continue
		}

		return authorsearch.Query{}, errors.New("only two names can be specified - last name and first name")
	}

	if queryStruct.LastName == "" {
		return authorsearch.Query{}, errors.New("last name has to be specified")
	}

	return queryStruct, nil
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
