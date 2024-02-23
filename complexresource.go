package main

import "strings"

func (website resource) searchComplexResource(query query) resource {
	// check history

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

	// possibly update history

	return website
}
