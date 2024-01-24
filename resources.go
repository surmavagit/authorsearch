package main

import "github.com/surmavagit/authorsearchcli/authorsearch"

var marxists = authorsearch.Resource{
	Name:       "marxists",
	BaseURL:    "https://www.marxists.org/",
	QueryURL:   "admin/js/data/authors.json",
	DataFormat: "json",
	URLFilter:  "index.htm",
}

var resources = []authorsearch.Resource{marxists}
