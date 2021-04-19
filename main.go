package main

import (
	"log"
	"net"
	"os"

	"github.com/asaskevich/govalidator"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
	"github.com/pkg/errors"
)

func main() {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	domains := []string{"perduuu.com", "maps.google.com", "perdu.com", "facebook.com", "jdkgjigikjb.net"}

	t.AppendHeader(table.Row{"#", "Domain name", "Available"})
	for idx, domain := range domains {
		exists, err := domainExists(domain)
		if err != nil {
			log.Println(err.Error())
		}
		if exists {
			t.AppendRows([]table.Row{{text.FgHiBlue.Sprintf("%v", idx), text.FgHiBlue.Sprintf("%v", domain), text.FgHiBlue.Sprintf("%v", !exists)}})
		} else {
			t.AppendRows([]table.Row{{text.FgHiGreen.Sprintf("%v", idx), text.FgHiGreen.Sprintf("%v", domain), text.FgHiGreen.Sprintf("%v", !exists)}})
		}
	}

	t.AppendSeparator()

	t.Render()
}

func domainExists(domain string) (bool, error) {
	if govalidator.IsDNSName(domain) {
		addrs, err := net.LookupHost(domain)
		if err != nil {
			return false, errors.Wrap(err, "LookupHost failed")
		}

		if len(addrs) == 0 {
			return false, errors.Errorf("Empty result")
		}
		return true, nil
	} else {
		return false, errors.Errorf("Not a valid domain name")
	}
}
