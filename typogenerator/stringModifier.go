package typogenerator

import (
	"strings"
)

func removeIndex(s string, i int) string {
	a := s[:i] + s[i+1:]
	return a
}

func SplitDomain(domain string) (string, string) {
	lastIdx := strings.LastIndex(domain, ".")
	domainTmp := domain[:lastIdx]
	topLevelDomain := domain[lastIdx+1:]
	return domainTmp, topLevelDomain
}

// skipLetter - removes letter from the original domain (example » xample, eample)
func SkipLetter(domain string) []string {
	var res []string
	for i := 0; i < len(domain); i++ {
		res = append(res, removeIndex(domain, i))
	}
	return res
}

// Insert Letter (example » erxample, edxample)
func InsertLetter(domain string) []string {
	var res []string
	for i := 0; i < len(domain); i++ {
		for j := 97; j < 123; j++ {
			tmp := domain[:i] + string(rune(j)) + domain[i:]
			res = append(res, tmp)
		}
	}
	return res
}

// Double Letter (example » eexample, exxample)
func DoubleLetter(domain string) []string {
	var res []string
	for i := 0; i < len(domain); i++ {
		tmp := domain[:i] + string(domain[i]) + domain[i:]
		res = append(res, tmp)
	}
	return res
}

// Strip Dashes These typos are created by omitting a dash from the domainname. For example, www.domain-name.com becomes www.domainname.com
func StripDashes(domain string) string {
	if strings.ContainsRune(domain, '-') {
		return strings.ReplaceAll(domain, "-", "")
	}
	return ""
}
