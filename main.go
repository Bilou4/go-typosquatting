package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"strings"

	"github.com/asaskevich/govalidator"
	"github.com/bilou4/go-typosquatting/typogenerator"
	"github.com/jedib0t/go-pretty/table"
	"github.com/jedib0t/go-pretty/text"
	"github.com/pkg/errors"
)

var (
	domain         string
	verbose        bool
	strategies     string
	strategiesList []string
	domains        []string
	Usage          = func(msg string) {
		fmt.Fprintf(os.Stderr, "[Usage] of %s:\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "[ERROR] => %s\n", msg)
		fmt.Fprintf(os.Stderr, "[Example] => %s -domain perdu.com -strategies skip,double -verbose", os.Args[0])
		flag.PrintDefaults()
		os.Exit(1)
	}
	ALL_STRATEGIES = []string{STRATEGY_SKIP, STRATEGY_INSERT, STRATEGY_DOUBLE, STRATEGY_STRIP_DASHES, STRATEGY_WRONG, STRATEGY_SWAP, STRATEGY_SWAP_VOWEL, STRATEGY_DOT, STRATEGY_HOMOGLYPHS, STRATEGY_TOP_DOMAIN}
)

const (
	STRATEGY_SKIP         = "skip"
	STRATEGY_INSERT       = "insert"
	STRATEGY_DOUBLE       = "double"
	STRATEGY_STRIP_DASHES = "strip-dashes"
	STRATEGY_WRONG        = "wrong"
	STRATEGY_SWAP         = "swap"
	STRATEGY_SWAP_VOWEL   = "swap-vowel"
	STRATEGY_DOT          = "dot"
	STRATEGY_HOMOGLYPHS   = "homoglyphs"
	STRATEGY_TOP_DOMAIN   = "tld"
)

func init() {
	flag.StringVar(&domain, "domain", "", "The domain name you want to usurp")
	flag.BoolVar(&verbose, "verbose", false, "To see log print.")
	helpStrategies := "Comma-separated list of strategies you want to use {"

	for idx, name := range ALL_STRATEGIES {
		if idx == len(ALL_STRATEGIES)-1 {
			helpStrategies += name + "}."
		} else {
			helpStrategies += name + "|"
		}
	}
	flag.StringVar(&strategies, "strategies", STRATEGY_DOUBLE, helpStrategies)

	flag.Parse()
	if domain == "" {
		Usage("Domain must not be empty")
	}
	strategiesList = strings.Split(strategies, ",")
	for _, strategy := range strategiesList {
		if !typogenerator.StringInSlice(strategy, ALL_STRATEGIES) {
			Usage("Not a valid strategy name -> " + strategy)
		}
	}
}

func main() {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	domainTmp, topLevelDomain := typogenerator.SplitDomain(domain)

	for _, strategy := range strategiesList {
		switch strategy {
		case STRATEGY_DOUBLE:
			domains = append(domains, typogenerator.DoubleLetter(domainTmp, topLevelDomain)...)
		case STRATEGY_SKIP:
			domains = append(domains, typogenerator.SkipLetter(domainTmp, topLevelDomain)...)
		case STRATEGY_INSERT:
			domains = append(domains, typogenerator.InsertLetter(domainTmp, topLevelDomain)...)
		case STRATEGY_STRIP_DASHES:
			domains = append(domains, typogenerator.StripDashes(domainTmp, topLevelDomain))
		case STRATEGY_WRONG:
			domains = append(domains, typogenerator.WrongLetter(domainTmp, topLevelDomain)...)
		case STRATEGY_SWAP:
			domains = append(domains, typogenerator.SwapLetter(domainTmp, topLevelDomain)...)
		case STRATEGY_SWAP_VOWEL:
			domains = append(domains, typogenerator.SwapVowel(domainTmp, topLevelDomain)...)
		case STRATEGY_DOT:
			domains = append(domains, typogenerator.MissingDot(domainTmp, topLevelDomain)...)
		case STRATEGY_HOMOGLYPHS:
			domains = append(domains, typogenerator.ReplaceByHomoglyphs(domainTmp, topLevelDomain)...)
		case STRATEGY_TOP_DOMAIN:
			domains = append(domains, typogenerator.ChangeTopDomain(domainTmp)...)
		default:
		}
	}

	t.AppendHeader(table.Row{"#", "Domain name", "Available"})
	for idx, domain := range domains {
		exists, err := domainExists(domain)
		if err != nil && verbose {
			log.Println(err.Error())
		}
		if exists {
			t.AppendRows([]table.Row{{text.FgHiBlue.Sprintf("%v", idx), text.FgHiBlue.Sprintf("%v", domain), text.FgHiBlue.Sprintf("%v", !exists)}})
		} else {
			t.AppendRows([]table.Row{{text.FgHiGreen.Sprintf("%v", idx), text.FgHiGreen.Sprintf("%v", domain), text.FgHiGreen.Sprintf("%v", !exists)}})
		}
	}

	// t.AppendSeparator()

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
		// it is consider as existing (true) otherwise, it will be print as available in the final tab.
		return true, errors.Errorf("Not a valid domain name %s", domain)
	}
}
