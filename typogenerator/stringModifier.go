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
