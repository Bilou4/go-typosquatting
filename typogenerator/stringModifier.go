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

// Wrong Letter - replaces a letter by another one (example » rxample, dxample)
func WrongLetter(domain string) []string {
	var res []string
	for i := 0; i < len(domain); i++ {
		for j := 97; j < 123; j++ {
			res = append(res, strings.Replace(domain, string(domain[i]), string(rune(j)), 1))
		}
	}
	return res
}

// SwapLetter - exchange two letters (example » xeample, eaxmple)
func SwapLetter(domain string) []string {
	var res []string
	for i := 0; i < len(domain)-1; i++ {
		tmp := []rune(domain)
		tmp[i], tmp[i+1] = tmp[i+1], tmp[i]
		res = append(res, string(tmp))
	}
	return res
}

// Vowel Swapping Swap vowels within the domain name except for the first letter. For example, www.google.com becomes www.gaagle.com.
func SwapVowel(domain string) []string {
	var res []string
	vowels := []string{"a", "e", "i", "o", "u", "y"}
	for i := 0; i < len(domain); i++ {
		for _, v := range vowels {
			if stringInSlice(string(domain[i]), vowels) {
				if string(domain[i]) != v {
					tmp := string(domain[:i]) + v + string(domain[i+1:])
					res = append(res, tmp)
				}
			}
		}
	}
	return res
}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

// Missing Dot These typos are created by omitting a dot from the domainname. For example, wwwgoogle.com and www.googlecom
func MissingDot(domain string) []string {
	var res []string
	nbDot := strings.Count(domain, ".")
	for i := 0; i < nbDot; i++ {
		res = append(res, strings.Replace(domain, ".", "", 1))
	}
	return res
}
