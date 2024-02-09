package authorsearch

import (
	"encoding/json"
	"errors"
	"strings"
)

// parseCache turns a byte stream from a cache file into a slice of authorData structs.
func (website Resource) parseCache(file []byte) ([]authorData, error) {
	var sliceOfData []authorData
	if website.DataFormat == "json" {
		err := json.Unmarshal(file, &sliceOfData)
		return sliceOfData, err
	}

	if website.DataFormat == "html" {
		root := parseHTML(file)
		if root == nil {
			return []authorData{}, errors.New("can't parse html")
		}
		links := getLinkElements(root)

		for _, l := range links {
			data := authorData{}
			data.AuthorURL = getHrefAttr(l)
			content := getTextContent(l.FirstChild)
			data.Description = strings.TrimSpace(strings.ReplaceAll(content, "\n", " "))
			sliceOfData = append(sliceOfData, data)
		}
		return sliceOfData, nil
	}

	return []authorData{}, errors.New("unknown resource data format")
}

// validData checks the validity of the authorData struct by applying a set of tests.
// The filter string in the arguments has to be a substring of the data.AuthorURL
// in order for the struct to be valid.
func validData(data authorData, filter string) bool {
	if data.AuthorURL == "" || data.Description == "" {
		return false
	}
	// get rid of links to sections of the same html page
	if strings.HasPrefix(data.AuthorURL, "#") {
		return false
	}
	if !strings.Contains(data.AuthorURL, filter) {
		return false
	}
	return true
}

// filterAndDedupe takes a slice of authorData structs and returns it after
// throwing out all invalid and duplicate structs. The filter string
// in the arguments has to be a substring of the data.AuthorURL in order
// for the struct to be valid.
func filterAndDedupe(data []authorData, filter string) []authorData {
	uniqueData := []authorData{}
	dataMap := map[string]bool{}
	separator := ":::"

	for _, d := range data {
		if validData(d, filter) {
			dataString := d.Description + separator + d.AuthorURL
			dataMap[dataString] = false
		}
	}

	for i := range dataMap {
		desc, url, ok := strings.Cut(i, separator)
		if !ok {
			continue
		}
		uniqueData = append(uniqueData, authorData{Description: desc, AuthorURL: url})
	}
	return uniqueData
}
