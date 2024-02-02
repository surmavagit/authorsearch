package main

import "github.com/surmavagit/authorsearchcli/authorsearch"

var marxists = authorsearch.Resource{
	Name:       "marxists",
	BaseURL:    "https://www.marxists.org/",
	QueryURL:   "admin/js/data/authors.json",
	DataFormat: "json",
	URLFilter:  "index.htm",
}

var mcmaster = authorsearch.Resource{
	Name:       "mcmaster",
	BaseURL:    "https://socialsciences.mcmaster.ca/econ/ugcm/3ll3/",
	QueryURL:   "authors.html",
	DataFormat: "html",
	URLFilter:  "",
}

var resources = []authorsearch.Resource{marxists, mcmaster}
