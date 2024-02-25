package main

import (
	"errors"
	"os"
	"regexp"
	"sync"

	flag "github.com/spf13/pflag"
)

type result struct {
	Name     string
	BaseURL  string
	Data     []authorData
	Error    error
	CacheErr error
}

func main() {
	verbose := flag.BoolP("verbose", "v", false, "include 'not found' results")
	nonum := flag.BoolP("no-numbers", "n", false, "do not print the total number of results per resource")
	nodesc := flag.BoolP("no-description", "d", false, "do not print result description")
	flag.Parse()

	searchQuery, err := checkInput(flag.Args())
	if err != nil {
		os.Stderr.WriteString(err.Error())
		os.Exit(1)
	}

	cacheDirectory := "cache"
	err = createDirIfNotExist(cacheDirectory)
	if err != nil {
		os.Stderr.WriteString(err.Error())
		os.Exit(1)
	}

	dataChan := make(chan result)
	wg := sync.WaitGroup{}
	wg.Add(len(resources))
	for _, r := range resources {
		resource := r
		go func() {
			data, cacheErr, err := resource.searchResource(searchQuery, cacheDirectory)
			dataChan <- result{Name: resource.Name, BaseURL: resource.BaseURL, Data: data, Error: err, CacheErr: cacheErr}
		}()
	}
	go func() {
		for r := range dataChan {
			printResults(r, *verbose, *nonum, *nodesc)
			wg.Done()
		}
	}()
	wg.Wait()
}

type query struct {
	LastName  string
	FirstName string
	Year      string
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
