package main

import (
	"encoding/json"
	"errors"
	"strings"
)

// readResource turns a byte stream from a response body into a slice of data.
func (website resource) readResource(file []byte) ([]authorData, error) {
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
			if !validURL(data.AuthorURL, website.URLFilter) {
				continue
			}

			data.AuthorURL = strings.TrimLeft(data.AuthorURL, "/")
			if data.AuthorURL == "" {
				continue
			}

			if website.DescInParent {
				if l.Parent == nil {
					return []authorData{}, errors.New("no parent html element")
				}
				data.Description = getTextContent(l.Parent.FirstChild)
			} else {
				data.Description = getTextContent(l.FirstChild)
			}

			oneLine := strings.ReplaceAll(data.Description, "\n", " ")
			data.Description = strings.Join(strings.Fields(oneLine), " ")
			if data.Description == "" {
				continue
			}

			sliceOfData = append(sliceOfData, data)
		}
		return sliceOfData, nil
	}

	return []authorData{}, errors.New("unknown resource data format")
}

// validURL checks the validity of the url by applying a set of tests.
// The filter string in the arguments has to be a substring of the url
// in order for it to be valid.
func validURL(url string, filter string) bool {
	if url == "" {
		return false
	}
	// get rid of links to sections of the same html page
	if strings.HasPrefix(url, "#") {
		return false
	}
	if !strings.Contains(url, filter) {
		return false
	}
	return true
}

// filterAndDedupe takes a slice of authorData structs and returns it after
// throwing out all invalid and duplicate structs.
func (website resource) dedupe(data []authorData) []authorData {
	uniqueData := []authorData{}
	dataMap := map[string]bool{}
	separator := "%%"

	for _, d := range data {
		dataString := d.AuthorURL + separator + d.Description
		dataMap[dataString] = false
	}

	for i := range dataMap {
		url, desc, ok := strings.Cut(i, separator)
		if !ok {
			continue
		}
		url = website.BaseURL + url
		uniqueData = append(uniqueData, authorData{Description: desc, AuthorURL: url})
	}
	return uniqueData
}
