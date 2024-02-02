package authorsearch

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestParseCache(t *testing.T) {
	website := Resource{DataFormat: "json"}
	expectedStruct := authorData{Description: "author name", AuthorURL: "www.example.com"}
	jsonToTest, err := json.Marshal([]authorData{expectedStruct})
	if err != nil {
		t.Fatal("Testing error: can't marshal json")
	}

	output, err := website.parseCache([]byte(jsonToTest))
	if err != nil || len(output) != 1 || output[0].AuthorURL != "www.example.com" || output[0].Description != "author name" {
		t.FailNow()
	}
}

func TestGetLinkElements(t *testing.T) {
	linkOne := authorData{Description: "test", AuthorURL: "www.example.com/1"}
	linkTwo := authorData{Description: "number one", AuthorURL: "www.example.com/2"}
	linkThree := authorData{Description: "and more text", AuthorURL: "www.example.com/anotherone"}
	expectedLinks := []authorData{linkOne, linkTwo, linkThree}
	htmlToCheck := fmt.Sprintf(`<div>
<p>this is a <a href="%s">%s</a></p>
<ul>
<li><a href="%s">%s</a></li>
<li>number two</li>
</ul>
    <a href="%s">%s</a>
</div>
`, linkOne.AuthorURL, linkOne.Description, linkTwo.AuthorURL, linkTwo.Description, linkThree.AuthorURL, linkThree.Description)

	root := parseHTML([]byte(htmlToCheck))
	if root == nil {
		t.Fatal("cannot find html root element")
	}

	linkSlice := getLinkElements(root)
	if len(linkSlice) != 3 {
		t.Fatalf("Expected 3 links, found %d", len(linkSlice))
	}
	for i, linkNode := range linkSlice {
		authorDataStruct := getDataFromHTML(linkNode)
		if !compareData(authorDataStruct, expectedLinks[i]) {
			t.Errorf("Expected Description: %s, URL: %s; Found Description: %s, URL: %s", expectedLinks[i].Description, expectedLinks[i].AuthorURL, authorDataStruct.Description, authorDataStruct.AuthorURL)
		}
	}
}

func compareData(one authorData, two authorData) bool {
	if one.AuthorURL == two.AuthorURL && one.Description == two.Description {
		return true
	}
	return false
}
