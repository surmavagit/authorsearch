package main

import (
	"encoding/json"
	"errors"
	"strings"

	"golang.org/x/net/html"
)

// readResource turns a byte stream from a response body into a slice of data.
func (website resource) readResource(file []byte) ([]authorData, error) {
	var sliceOfData []authorData
	if website.DataFormat == "json" {
		err := json.Unmarshal(file, &sliceOfData)
		return sliceOfData, err
	}

	if website.DataFormat == "html" {
		root, err := parseHTML(file)
		if err != nil {
			return nil, err
		}
		links := getLinkElements(root)

		for _, l := range links {
			data, ok := website.getDataFromLink(l)
			if !ok {
				continue
			}
			sliceOfData = append(sliceOfData, data)
		}
		return sliceOfData, nil
	}

	return nil, errors.New("unknown resource data format")
}

func (website resource) getDataFromLink(l *html.Node) (authorData, bool) {
	data := authorData{}

	href := getHrefAttr(l)
	if !validURL(href, website.URLFilter) {
		return authorData{}, false
	}
	data.AuthorURL = strings.TrimLeft(href, "/")
	if data.AuthorURL == "" {
		return authorData{}, false
	}

	desc := ""
	if website.DescInParent {
		if l.Parent == nil {
			return authorData{}, false
		}
		desc = getTextContent(l.Parent.FirstChild)
	} else {
		desc = getTextContent(l.FirstChild)
	}
	oneLine := strings.ReplaceAll(desc, "\n", " ")
	data.Description = strings.Join(strings.Fields(oneLine), " ")
	if data.Description == "" {
		return authorData{}, false
	}

	return data, true
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

// dedupe takes a slice of authorData structs and returns it after
// throwing out all invalid and duplicate structs.
func dedupe(data []authorData) []authorData {
	uniqueData := []authorData{}
	dataMap := map[string]string{}

	for _, d := range data {
		dataMap[d.AuthorURL] = d.Description
	}

	for u, d := range dataMap {
		uniqueData = append(uniqueData, authorData{Description: d, AuthorURL: u})
	}
	return uniqueData
}
