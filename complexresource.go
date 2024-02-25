package main

import (
	"strings"
)

type records map[string][]authorData

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
