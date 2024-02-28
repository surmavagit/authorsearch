# authorsearch
Search for an author across several library websites.

- quickly search multiple resources from the command line
- each resource is handled concurrently
- results are cached to reduce duplicate requests
- several formatting options
- errors go to standard error stream to allow filtering them out

## Quickstart
Here you can [install Go](https://go.dev/doc/install), if you don't have it already.

To install locally, first run ``git clone github.com/surmavagit/authorsearch`` to clone this repository.

Then run ``go build`` to make an executable file and then run it.

As an alternative, you can first install Go and then run ``go install github.com/surmavagit/authorsearch`` to install it globally.

## Usage
- Run ``./authorsearch Lastname`` to search for an author with this name:

e.g. ``./authorsearch Smith`` will give you results for authors called Smith.

- You can also additionaly specify author's first name or year of birth or year of death to narrow down the results:

e.g. ``./authorsearch Smith Adam`` or ``./authorsearch Smith 1790`` or ``./authorsearch Smith Adam 1723`` are all valid requests.

- There are several optional flags to modify the presentation of results, to check them run ``./authorsearch --help`` or ``-h``



## Featured Libraries

### General Purpose 
1. [OpenLibrary](openlibrary.org)

a project of the Internet Archive, hosts thousands of books on all topics in different formats (picture and text)

2. [Project Gutenberg](gutenberg.org)

hosts thousands of books on all topics in html or other text formats

3. [OnlineBooksPage](onlinebooks.library.upenn.edu)

while it doesn't host any books itself, the OnlineBooksPage provides outstanding information on where one can read free public domain texts on all topics

### Libraries on the topics of history and political economy:
4. [Tokyo Keizai University Institutional Repository](repository.tku.ac.jp) - multiple texts (page scan pictures)
5. [Marxists Internet Archive](marxists.org) - multiple texts not only by Marx, but also by many other authors
6. [History of Economic Thought website](hetwebsite.net/het/) - information about authors and multiple links to their texts
7. [Mcmaster University Archive for the History of Economic Thought](socialsciences.mcmaster.ca/econ/ugcm/3ll3/) - multiple texts in txt or pdf formats
8. [Site Paulette Taieb](taieb.net) - texts in both english and french (originals and translations)

