package main

type resource struct {
	Name         string
	Complex      bool // false, if the website contains a single list of all authors, true otherwise
	BaseURL      string
	QueryURL     string
	DataFormat   string // json or html
	DescInParent bool   // Location of the full description in HTML document relative to the <a> tag
	URLFilter    string // Valid URLs contain this string
	FirstName    bool
	Year         bool
}

var resources = map[string]resource{
	"marxists": {
		Name:         "marxists",
		Complex:      false,
		BaseURL:      "https://www.marxists.org/",
		QueryURL:     "admin/js/data/authors.json",
		DataFormat:   "json",
		DescInParent: false,
		URLFilter:    "",
		FirstName:    true,
		Year:         false,
	},

	"mcmaster": {
		Name:         "mcmaster",
		Complex:      false,
		BaseURL:      "https://socialsciences.mcmaster.ca/econ/ugcm/3ll3/",
		QueryURL:     "authors.html",
		DataFormat:   "html",
		DescInParent: false,
		URLFilter:    "",
		FirstName:    true,
		Year:         false,
	},

	"hetwebsite": {
		Name:         "hetwebsite",
		Complex:      false,
		BaseURL:      "https://www.hetwebsite.net/het/",
		QueryURL:     "alphabet.htm",
		DataFormat:   "html",
		DescInParent: true,
		URLFilter:    "profiles",
		FirstName:    true,
		Year:         true,
	},

	"taieb": {
		Name:         "taieb",
		Complex:      false,
		BaseURL:      "https://www.taieb.net/",
		QueryURL:     "menu/index.html",
		DataFormat:   "html",
		DescInParent: false,
		URLFilter:    "",
		FirstName:    false,
		Year:         false,
	},

	"gutenberg": {
		Name:         "gutenberg",
		Complex:      true,
		BaseURL:      "https://www.gutenberg.org/",
		QueryURL:     "ebooks/authors/search/?query=",
		DataFormat:   "html",
		DescInParent: false,
		URLFilter:    "ebooks/author/",
		FirstName:    true,
		Year:         true,
	},

	"openlib": {
		Name:         "openlib",
		Complex:      true,
		BaseURL:      "https://openlibrary.org/",
		QueryURL:     "search/authors?q=",
		DataFormat:   "html",
		DescInParent: true,
		URLFilter:    "authors/",
		FirstName:    true,
		Year:         false,
	},
}
