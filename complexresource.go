package main

import (
	"strings"
)

type records map[string][]authorData

func (website resource) searchComplexResource(query query, cacheDir string) (data []authorData, err error) {
	// normalise query - only for complex resources
	query.LastName = strings.ToLower(query.LastName)
	query.FirstName = strings.ToLower(query.FirstName)

	// check history
	histFile := cacheDir + "/" + website.Name + "_" + query.LastName + ".json"
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

		found, data := searchInHistory(history, query)
		if found {
			return data, nil
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
	body, err := requestURL(fullURL)
	if err != nil {
		return []authorData{}, err
	}

	rawData, err := website.readResource(body)
	if err != nil {
		return []authorData{}, err
	}

	filteredData := []authorData{}
	// filter by year if the resource doesn't filter itself
	if !website.Year {
		for _, d := range rawData {
			if strings.Contains(d.Description, query.Year) {
				filteredData = append(filteredData, d)
			}
		}
	}

	// update history
	queryString := strings.Join(queryTerms, " ")
	history[queryString] = filteredData

	return filteredData, writeFileJSON(histFile, history)
}

func searchInHistory(history records, query query) (bool, []authorData) {
	last, ok := history[query.LastName]
	if ok {
		if last == nil || query.FirstName == "" && query.Year == "" {
			return true, last
		}
	}

	if query.FirstName != "" {
		first, ok := history[query.LastName+" "+query.FirstName]
		if ok {
			if first == nil || query.Year == "" {
				return true, first
			}
		}
	}

	if query.Year != "" {
		year, ok := history[query.LastName+" "+query.Year]
		if ok {
			if year == nil || query.FirstName == "" {
				return true, year
			}
		}
	}

	if query.FirstName != "" && query.Year != "" {
		full, ok := history[query.LastName+" "+query.FirstName+" "+query.Year]
		if ok {
			return true, full
		}
	}

	return false, nil
}
