package main

var marxists = resource{
	Name:         "marxists",
	BaseURL:      "https://www.marxists.org/",
	QueryURL:     "admin/js/data/authors.json",
	DataFormat:   "json",
	DescInParent: false,
	URLFilter:    "",
	FirstName:    true,
	Year:         false,
}

var mcmaster = resource{
	Name:         "mcmaster",
	BaseURL:      "https://socialsciences.mcmaster.ca/econ/ugcm/3ll3/",
	QueryURL:     "authors.html",
	DataFormat:   "html",
	DescInParent: false,
	URLFilter:    "",
	FirstName:    true,
	Year:         false,
}

var hetwebsite = resource{
	Name:         "hetwebsite",
	BaseURL:      "https://www.hetwebsite.net/het/",
	QueryURL:     "alphabet.htm",
	DataFormat:   "html",
	DescInParent: true,
	URLFilter:    "profiles",
	FirstName:    true,
	Year:         true,
}

var taieb = resource{
	Name:         "taieb",
	BaseURL:      "https://www.taieb.net/",
	QueryURL:     "menu/index.html",
	DataFormat:   "html",
	DescInParent: false,
	URLFilter:    "",
	FirstName:    false,
	Year:         false,
}

var resources = []resource{marxists, mcmaster, hetwebsite, taieb}
