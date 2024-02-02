package authorsearch

import (
	"bytes"
	"encoding/json"
	"errors"
	"strings"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
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

		// get rid of duplicates
		links := getLinkElements(root)
		dataMap := map[string]bool{}
		for _, l := range links {
			dataString := getDataStringFromHTML(l)
			if dataString != "" {
				dataMap[dataString] = false
			}
		}

		for i := range dataMap {
			dataStruct, err := getDataStructFromString(i)
			if err != nil {
				return []authorData{}, err
			}
			sliceOfData = append(sliceOfData, dataStruct)
		}

		return sliceOfData, nil
	}

	return []authorData{}, errors.New("unknown resource data format")
}

// parseHTML returns a pointer to the root node of html or, if parsing fails, it returns nil
func parseHTML(file []byte) *html.Node {
	reader := bytes.NewReader(file)
	rootNode, _ := html.Parse(reader)
	return rootNode
}

// getLinkElements takes a pointer to an html node, recursively goes through all its children and siblings
// and returns a slice of pointers to all 'a' html nodes.
func getLinkElements(element *html.Node) []*html.Node {
	var linkElements []*html.Node

	if element.DataAtom == atom.A {
		linkElements = append(linkElements, element)
	} else if element.FirstChild != nil {
		linkElements = append(linkElements, getLinkElements(element.FirstChild)...)
	}

	if element.NextSibling != nil {
		linkElements = append(linkElements, getLinkElements(element.NextSibling)...)
	}

	return linkElements
}

// getDataStringFromHTML takes a pointer to an html node, representing an author,
// and turns it into a string, that combines description and url. Returns an empty string
// if either description or url is missing.
func getDataStringFromHTML(link *html.Node) string {
	if link.FirstChild == nil {
		return ""
	}
	description := strings.ReplaceAll(link.FirstChild.Data, "\n", " ")

	urlLink := ""
	for _, a := range link.Attr {
		if a.Key == "href" && a.Val == "" {
			return ""
		} else if a.Key == "href" {
			urlLink = a.Val
			break
		}
	}
	return description + ":::" + urlLink
}

func getDataStructFromString(data string) (authorData, error) {
	desc, url, ok := strings.Cut(data, ":::")
	if !ok {
		return authorData{}, errors.New("wrong Data string: " + data)
	}
	return authorData{Description: desc, AuthorURL: url}, nil
}

func (website Resource) filterData(rawData []authorData) ([]authorData, error) {
	if website.URLFilter == "" {
		return rawData, nil
	}
	var filteredData []authorData
	for _, d := range rawData {
		if strings.Contains(d.AuthorURL, website.URLFilter) {
			filteredData = append(filteredData, d)
		}
	}
	return filteredData, nil
}
