package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"sync"

	"github.com/asaskevich/govalidator"
	"github.com/bilou4/go-typosquatting/typogenerator"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
	"github.com/pkg/errors"
)

var (
	cs             ConcurrentSlice
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
	wg             sync.WaitGroup
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

// ConcurrentSlice allows to access the domains string slice with multiples goroutines
type ConcurrentSlice struct {
	sync.RWMutex
	domains []string
}

// Appends an item to the concurrent slice
func (cs *ConcurrentSlice) Append(items []string) {
	cs.Lock()
	defer cs.Unlock()
	cs.domains = append(cs.domains, items...)
}

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

	// we need to wait X goroutines (X is the number of strategies defined by the user)
	wg.Add(len(strategiesList))

	for _, strategy := range strategiesList {
		switch strategy {
		case STRATEGY_DOUBLE:
			go func(wg *sync.WaitGroup, cs *ConcurrentSlice) {
				defer wg.Done()
				cs.Append(typogenerator.DoubleLetter(domainTmp, topLevelDomain))
			}(&wg, &cs)
		case STRATEGY_SKIP:
			go func(wg *sync.WaitGroup, cs *ConcurrentSlice) {
				defer wg.Done()
				cs.Append(typogenerator.SkipLetter(domainTmp, topLevelDomain))
			}(&wg, &cs)
		case STRATEGY_INSERT:
			go func(wg *sync.WaitGroup, cs *ConcurrentSlice) {
				defer wg.Done()
				cs.Append(typogenerator.InsertLetter(domainTmp, topLevelDomain))
			}(&wg, &cs)
		case STRATEGY_STRIP_DASHES:
			go func(wg *sync.WaitGroup, cs *ConcurrentSlice) {
				defer wg.Done()
				c := typogenerator.StripDashes(domainTmp, topLevelDomain)
				var s []string = []string{c}
				cs.Append(s)
			}(&wg, &cs)
		case STRATEGY_WRONG:
			go func(wg *sync.WaitGroup, cs *ConcurrentSlice) {
				defer wg.Done()
				cs.Append(typogenerator.WrongLetter(domainTmp, topLevelDomain))
			}(&wg, &cs)
		case STRATEGY_SWAP:
			go func(wg *sync.WaitGroup, cs *ConcurrentSlice) {
				defer wg.Done()
				cs.Append(typogenerator.SwapLetter(domainTmp, topLevelDomain))
			}(&wg, &cs)
		case STRATEGY_SWAP_VOWEL:
			go func(wg *sync.WaitGroup, cs *ConcurrentSlice) {
				defer wg.Done()
				cs.Append(typogenerator.SwapVowel(domainTmp, topLevelDomain))
			}(&wg, &cs)
		case STRATEGY_DOT:
			go func(wg *sync.WaitGroup, cs *ConcurrentSlice) {
				defer wg.Done()
				cs.Append(typogenerator.MissingDot(domainTmp, topLevelDomain))
			}(&wg, &cs)
		case STRATEGY_HOMOGLYPHS:
			go func(wg *sync.WaitGroup, cs *ConcurrentSlice) {
				defer wg.Done()
				cs.Append(typogenerator.ReplaceByHomoglyphs(domainTmp, topLevelDomain))
			}(&wg, &cs)
		case STRATEGY_TOP_DOMAIN:
			go func(wg *sync.WaitGroup, cs *ConcurrentSlice) {
				defer wg.Done()
				cs.Append(typogenerator.ChangeTopDomain(domainTmp))
			}(&wg, &cs)
		default:
		}
	}
	wg.Wait()

	t.AppendHeader(table.Row{"#", "Domain name", "Available"})
	for idx, domain := range cs.domains {
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
		// it is consider as existing (true) otherwise, it will be print as available in the final tab.
		return true, errors.Errorf("Not a valid domain name %s", domain)
	}
}
