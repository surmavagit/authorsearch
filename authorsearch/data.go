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
			data := authorData{}
			data.AuthorURL = getHrefAttr(l)
			content := getTextContent(l.FirstChild)
			data.Description = strings.ReplaceAll(content, "\n", " ")
			sliceOfData = append(sliceOfData, data)
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

// getLinkElements recursively goes through all the children and siblings of the provided html node
// and returns a slice of pointers to all 'a' html nodes.
func getLinkElements(element *html.Node) []*html.Node {
	if element == nil {
		return []*html.Node{}
	}

	if element.DataAtom == atom.A {
		return append([]*html.Node{element}, getLinkElements(element.NextSibling)...)
	}

	return append(getLinkElements(element.FirstChild), getLinkElements(element.NextSibling)...)
}

// returns the value of the href attribute of the provided html node or empty string
func getHrefAttr(link *html.Node) string {
	for _, a := range link.Attr {
		if a.Key == "href" {
			return a.Val
		}
	}
	return ""
}

// getTextContent recursively goes through all the children and siblings of the provided html node
// and returns their combined text content
func getTextContent(element *html.Node) string {
	if element == nil {
		return ""
	}

	if element.Type == html.TextNode {
		return element.Data + getTextContent(element.NextSibling)
	}

	return getTextContent(element.FirstChild) + getTextContent(element.NextSibling)
}

func validData(data authorData, filter string) bool {
	if data.AuthorURL == "" || data.Description == "" {
		return false
	}
	if strings.HasPrefix(data.AuthorURL, "#") {
		return false
	}
	if !strings.Contains(data.AuthorURL, filter) {
		return false
	}
	return true
}

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
