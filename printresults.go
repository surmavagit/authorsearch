package main

import (
	"fmt"
	"os"
)

func printResults(r result, verbose bool, nonum bool, nodesc bool) {
	if r.CacheErr != nil {
		os.Stderr.WriteString(fmt.Sprintln(formatName(r.Name), r.CacheErr.Error()))
	}

	if r.Error != nil {
		os.Stderr.WriteString(fmt.Sprintln(formatName(r.Name), r.Error.Error()))
		return
	}

	numLinks := len(r.Data)
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

	for i, l := range r.Data {
		url := r.BaseURL + l.AuthorURL
		if nonum && nodesc {
			fmt.Println(formatName(r.Name), url)
			continue
		}

		if nonum {
			fmt.Println(formatName(r.Name), formatDesc(l.Description, nonum), url)
			continue
		}

		num := fmt.Sprintf("[%d of %d]", i+1, numLinks)
		if nodesc {
			fmt.Println(formatName(r.Name), formatNum(num), url)
			continue
		}
		fmt.Println(formatName(r.Name), formatNum(num), formatDesc(l.Description, nonum), url)
	}
}

func formatName(name string) string {
	return fmt.Sprintf("%-11s", name)
}

func formatNum(num string) string {
	return fmt.Sprintf("%-10s", num)
}

func formatDesc(desc string, nonum bool) string {
	if nonum {
		return fmt.Sprintf("%-46s", shorten(desc, 46))
	}
	return fmt.Sprintf("%-35s", shorten(desc, 35))
}

func shorten(s string, maximum int) string {
	if len(s) <= maximum {
		return s
	}

	return s[0:maximum-3] + "..."
}
