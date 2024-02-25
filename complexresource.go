package main

import (
	"strings"
)

type records map[string][]authorData

func (website resource) searchComplexResource(q query, histFile string) ([]authorData, error) {
	// check history
	noHistory, err := fileNotExist(histFile)
	if err != nil {
		return []authorData{}, err
	}

	// load history and search
	history := records{}
	if !noHistory {
		err := loadFileJSON(histFile, &history)
		if err != nil {
			return []authorData{}, err
		}

		found, data := searchInHistory(history, q)
		if found {
			return data, nil
		}
	}

	// query and analyze result
	rawData, err := website.getResource(q)
	if err != nil {
		return []authorData{}, err
	}

	filteredData := website.filterRelevant(rawData, q)

	// update history
	queryString := getQueryString(q)
	history[queryString] = filteredData

	return filteredData, writeFileJSON(histFile, history)
}

func getQueryString(q query) string {
	q.LastName = strings.ToLower(q.LastName)
	q.FirstName = strings.ToLower(q.FirstName)

	querySlice := []string{q.LastName}

	if q.FirstName != "" {
		querySlice = append(querySlice, q.FirstName)
	}
	if q.Year != "" {
		querySlice = append(querySlice, q.Year)
	}
	return strings.Join(querySlice, " ")
}

func searchInHistory(history records, q query) (bool, []authorData) {
	q.LastName = strings.ToLower(q.LastName)
	q.FirstName = strings.ToLower(q.FirstName)

	last, ok := history[q.LastName]
	if ok {
		if last == nil || q.FirstName == "" && q.Year == "" {
			return true, last
		}
	}

	if q.FirstName != "" {
		first, ok := history[q.LastName+" "+q.FirstName]
		if ok {
			if first == nil || q.Year == "" {
				return true, first
			}
		}
	}

	if q.Year != "" {
		year, ok := history[q.LastName+" "+q.Year]
		if ok {
			if year == nil || q.FirstName == "" {
				return true, year
			}
		}
	}

	if q.FirstName != "" && q.Year != "" {
		full, ok := history[q.LastName+" "+q.FirstName+" "+q.Year]
		if ok {
			return true, full
		}
	}

	return false, nil
}
