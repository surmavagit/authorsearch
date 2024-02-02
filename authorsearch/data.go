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

		links := getLinkElements(root)
		for _, l := range links {
			sliceOfData = append(sliceOfData, getDataFromHTML(l))
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

// getDataFromHTML takes a pointer to an html node, representing an author,
// and turns it into a authorData struct
func getDataFromHTML(link *html.Node) authorData {
	var author authorData
	for _, a := range link.Attr {
		if a.Key == "href" {
			author.AuthorURL = a.Val
			break
		}
	}
	author.Description = link.FirstChild.Data
	return author
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
