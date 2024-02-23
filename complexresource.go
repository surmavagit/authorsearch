package main

import (
	"os"
	"strings"
)

type records map[string][]authorData

func (website resource) searchComplexResource(query query, cacheDir string) resource {
	// normalise query - only for complex resources
	query.LastName = strings.ToLower(query.LastName)
	query.FirstName = strings.ToLower(query.FirstName)
	query.Year = strings.ToLower(query.Year)

	// check history
	histFile := cacheDir + "/" + website.Name
	noHistory, err := fileNotExist(histFile)
	if err != nil && !noHistory {
		website.Error = err
		return website
	}

	history := records{}
	if !noHistory {
		// load history
		err := loadFileJSON(histFile, &history)
		if err != nil {
			website.Error = err
			return website
		}

		// search in history

		last, ok := history[query.LastName]
		if ok {
			if last == nil || query.FirstName == "" && query.Year == "" {
				website.Results = last
				return website
			}
		}

		if query.FirstName != "" {
			first, ok := history[query.LastName+" "+query.FirstName]
			if ok {
				if first == nil || query.Year == "" {
					website.Results = first
					return website
				}
			}
		}

		if query.Year != "" {
			year, ok := history[query.LastName+" "+query.Year]
			if ok {
				if year == nil || query.FirstName == "" {
					website.Results = year
					return website
				}
			}
		}

		if query.FirstName != "" && query.Year != "" {
			full, ok := history[query.LastName+" "+query.FirstName+" "+query.Year]
			if ok {
				website.Results = full
				return website
			}
		}

	}

	// query and analyze result
	queryTerms := []string{query.LastName}
	if website.FirstName && query.FirstName != "" {
		queryTerms = append(queryTerms, query.FirstName)
	}
	if website.Year && query.Year != "" {
		queryTerms = append(queryTerms, query.Year)
	}
	formatQuery := strings.Join(queryTerms, "+")

	fullURL := website.BaseURL + website.QueryURL + formatQuery
	body, err := getResource(fullURL)
	if err != nil {
		website.Error = err
		return website
	}

	data, err := website.readResource(body)
	if err != nil {
		website.Error = err
		return website
	}

	website.Results = data

	// update history
	queryString := strings.Join(queryTerms, " ")
	history[queryString] = data
	err = writeFileJSON(histFile, history)
	if err != nil {
		os.Stderr.WriteString("WARNING: " + err.Error())
	}

	return website
}
