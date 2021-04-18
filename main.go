package main

import (
	"log"
	"os"
	"strings"

	"github.com/domainr/whois"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
	"github.com/pkg/errors"
)

func main() {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	domains := []string{"perduuu.com", "google.com", "perdu.com", "facebook.com", "jdkgjigikjb.net"}

	t.AppendHeader(table.Row{"#", "Domain name", "Available"})
	for idx, domain := range domains {
		exists, err := domainExists(domain)
		if err != nil {
			log.Fatalf(err.Error())
		}
		if exists {
			t.AppendRows([]table.Row{{text.FgHiBlue.Sprintf("%v", idx), text.FgHiBlue.Sprintf("%v", domain), text.FgHiBlue.Sprintf("%v", exists)}})
		} else {
			t.AppendRows([]table.Row{{text.FgHiGreen.Sprintf("%v", idx), text.FgHiGreen.Sprintf("%v", domain), text.FgHiGreen.Sprintf("%v", exists)}})
		}
	}

	t.AppendSeparator()

	t.Render()
}

func domainExists(domain string) (bool, error) {
	request, err := whois.NewRequest(domain)
	if err != nil {
		return false, errors.Wrap(err, "NewRequest failed")
	}
	response, err := whois.DefaultClient.Fetch(request)

	if err != nil {
		return false, errors.Wrap(err, "Fetch failed")
	}

	if strings.Contains(response.String(), "No match") {
		return false, nil
	}

	return true, nil
}
