package authorsearch

import (
	"encoding/json"
	"testing"
)

func TestParseCache(t *testing.T) {
	website := Resource{}
	expectedStruct := data{Description: "author name", AuthorURL: "www.example.com"}
	jsonToTest, err := json.Marshal([]data{expectedStruct})
	if err != nil {
		t.Fatal("Testing error: can't marshal json")
	}

	output, err := website.parseCache([]byte(jsonToTest))
	if err != nil || len(output) != 1 || output[0].AuthorURL != "www.example.com" || output[0].Description != "author name" {
		t.FailNow()
	}
}
