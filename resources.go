package main

import "github.com/surmavagit/authorsearchcli/authorsearch"

var marxists = authorsearch.Resource{
	Name:       "marxists",
	BaseURL:    "https://www.marxists.org",
	QueryURL:   "/admin/js/data/authors.json",
	DataFormat: "json",
	URLFilter:  "",
}

var mcmaster = authorsearch.Resource{
	Name:       "mcmaster",
	BaseURL:    "https://socialsciences.mcmaster.ca/econ/ugcm/3ll3",
	QueryURL:   "/authors.html",
	DataFormat: "html",
	URLFilter:  "",
}

var hetwebsite = authorsearch.Resource{
	Name:       "hetwebsite",
	BaseURL:    "https://www.hetwebsite.net/het",
	QueryURL:   "/alphabet.htm",
	DataFormat: "html",
	URLFilter:  "profiles",
}

var taieb = authorsearch.Resource{
	Name:       "taieb",
	BaseURL:    "https://www.taieb.net",
	QueryURL:   "/menu/index.html",
	DataFormat: "html",
	URLFilter:  "",
}

var resources = []authorsearch.Resource{marxists, mcmaster, hetwebsite, taieb}
