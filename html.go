package authorsearch

import (
	"bytes"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

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

// getHrefAttr returns the value of the href attribute of the provided html node or empty string
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
