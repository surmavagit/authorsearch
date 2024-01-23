package main

import "github.com/surmavagit/authorsearchcli/authorsearch"

var marxists = authorsearch.Resource{
	Name:      "marxists",
	BaseURL:   "https://www.marxists.org/",
	QueryURL:  "admin/js/data/authors.json",
	CacheFile: "cache/marxists.json",
}

var resources = []authorsearch.Resource{marxists}
