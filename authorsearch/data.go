package authorsearch

import (
	"encoding/json"
	"strings"
)

// parseCache turns a byte stream from a cache file into a slice of data structs.
func (website Resource) parseCache(file []byte) ([]data, error) {
	var rawData []data
	err := json.Unmarshal(file, &rawData)
	if err != nil {
		return []data{}, err
	}

	return rawData, nil
}

func (website Resource) filterData(rawData []data) ([]data, error) {
	if website.URLFilter == "" {
		return rawData, nil
	}
	var filteredData []data
	for _, d := range rawData {
		if strings.Contains(d.AuthorURL, website.URLFilter) {
			filteredData = append(filteredData, d)
		}
	}
	return filteredData, nil
}
