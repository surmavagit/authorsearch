package main

import (
	"fmt"
	"os"
)

func printResults(r resource, verbose bool, nonum bool, nodesc bool) {
	if r.Error != nil {
		os.Stderr.WriteString(fmt.Sprintln(formatName(r.Name), r.Error.Error()))
		return
	}

	numLinks := len(r.Results)
	if numLinks == 0 && !verbose {
		return
	}

	if numLinks == 0 {
		if nonum {
			fmt.Println(formatName(r.Name), "not found")
			return
		}
		num := "[0 of 0]"
		fmt.Println(formatName(r.Name), formatNum(num), "not found")
		return
	}

	for i, l := range r.Results {
		if nonum && nodesc {
			fmt.Println(formatName(r.Name), l.AuthorURL)
			continue
		}

		if nonum {
			fmt.Println(formatName(r.Name), formatDesc(l.Description, nonum), l.AuthorURL)
			continue
		}

		num := fmt.Sprintf("[%d of %d]", i+1, numLinks)
		if nodesc {
			fmt.Println(formatName(r.Name), formatNum(num), l.AuthorURL)
			continue
		}
		fmt.Println(formatName(r.Name), formatNum(num), formatDesc(l.Description, nonum), l.AuthorURL)
	}
}

func formatName(name string) string {
	return fmt.Sprintf("%-11s", name)
}

func formatNum(num string) string {
	return fmt.Sprintf("%-9s", num)
}

func formatDesc(desc string, nonum bool) string {
	if nonum {
		return fmt.Sprintf("%-45s", desc)
	}
	return fmt.Sprintf("%-35s", desc)
}
