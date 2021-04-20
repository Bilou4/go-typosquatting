package typogenerator

import (
	"strings"
)

func removeIndex(s string, i int) string {
	a := s[:i] + s[i+1:]
	return a
}

// skipLetter - removes letter from the original domain (example » xample, eample)
func SkipLetter(domain string) []string {
	var res []string
	splitDomain := strings.Split(domain, ".")
	for i := 0; i < len(splitDomain[0]); i++ {
		res = append(res, removeIndex(splitDomain[0], i))
	}
	return res
}

// Insert Letter (example » erxample, edxample)
func InsertLetter(domain string) []string {
	var res []string
	splitDomain := strings.Split(domain, ".")
	for i := 0; i < len(splitDomain[0]); i++ {
		for j := 97; j < 123; j++ {
			tmp := splitDomain[0][:i] + string(rune(j)) + splitDomain[0][i:]
			res = append(res, tmp)
		}
	}
	return res
}
