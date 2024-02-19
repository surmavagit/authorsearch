package main

import (
	"errors"
	"fmt"
	"os"
	"regexp"
	"sync"
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

	dataChan := make(chan resource)
	wg := sync.WaitGroup{}
	for _, r := range resources {
		resource := r
		wg.Add(1)
		go func() {
			dataChan <- resource.searchResource(searchQuery, cacheDirectory)
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

func printResults(r resource) {
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

func checkInput(input []string) (query, error) {
	if len(input) == 0 {
		return query{}, errors.New("no search query provided")
	}

	if len(input) > 3 {
		return query{}, errors.New("too many arguments")
	}

	queryStruct := query{}
	numeric, err := regexp.Compile(`\d`)
	if err != nil {
		return query{}, err
	}
	nonNumeric, err := regexp.Compile(`\D`)
	if err != nil {
		return query{}, err
	}

	for _, a := range input {
		num := numeric.MatchString(a)
		nonNum := nonNumeric.MatchString(a)

		if num && nonNum {
			return query{}, errors.New("invalid argument: numeric and nonnumeric characters")
		}

		if num && queryStruct.Year != "" {
			return query{}, errors.New("only one year can be specified")
		}

		if num {
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

		return query{}, errors.New("only two names can be specified - last name and first name")
	}

	if queryStruct.LastName == "" {
		return query{}, errors.New("last name has to be specified")
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
