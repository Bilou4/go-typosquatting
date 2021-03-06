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
func SkipLetter(domain, tld string) []string {
	var res []string
	for i := 0; i < len(domain); i++ {
		tmp := removeIndex(domain, i)
		res = append(res, concatDomainTopLevel(tmp, tld))
	}
	return res
}

// Insert Letter (example » erxample, edxample)
func InsertLetter(domain, tld string) []string {
	var res []string
	for i := 0; i < len(domain); i++ {
		for j := 97; j < 123; j++ {
			tmp := domain[:i] + string(rune(j)) + domain[i:]
			res = append(res, concatDomainTopLevel(tmp, tld))
		}
	}
	return res
}

// Double Letter (example » eexample, exxample)
func DoubleLetter(domain, tld string) []string {
	var res []string
	for i := 0; i < len(domain); i++ {
		tmp := domain[:i] + string(domain[i]) + domain[i:]
		res = append(res, concatDomainTopLevel(tmp, tld))
	}
	return res
}

// Strip Dashes These typos are created by omitting a dash from the domainname. For example, www.domain-name.com becomes www.domainname.com
func StripDashes(domain, tld string) string {
	if strings.ContainsRune(domain, '-') {
		return concatDomainTopLevel(strings.ReplaceAll(domain, "-", ""), tld)
	}
	return concatDomainTopLevel(domain, tld)
}

// Wrong Letter - replaces a letter by another one (example » rxample, dxample)
func WrongLetter(domain, tld string) []string {
	var res []string
	for i := 0; i < len(domain); i++ {
		for j := 97; j < 123; j++ {
			res = append(res, concatDomainTopLevel(strings.Replace(domain, string(domain[i]), string(rune(j)), 1), tld))
		}
	}
	return res
}

// SwapLetter - exchange two letters (example » xeample, eaxmple)
func SwapLetter(domain, tld string) []string {
	var res []string
	for i := 0; i < len(domain)-1; i++ {
		tmp := []rune(domain)
		tmp[i], tmp[i+1] = tmp[i+1], tmp[i]
		res = append(res, concatDomainTopLevel(string(tmp), tld))
	}
	return res
}

// Vowel Swapping Swap vowels within the domain name except for the first letter. For example, www.google.com becomes www.gaagle.com.
func SwapVowel(domain, tld string) []string {
	var res []string
	vowels := []string{"a", "e", "i", "o", "u", "y"}
	for i := 0; i < len(domain); i++ {
		for _, v := range vowels {
			if StringInSlice(string(domain[i]), vowels) {
				if string(domain[i]) != v {
					tmp := string(domain[:i]) + v + string(domain[i+1:])
					res = append(res, concatDomainTopLevel(tmp, tld))
				}
			}
		}
	}
	return res
}

func StringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

// Missing Dot These typos are created by omitting a dot from the domainname. For example, wwwgoogle.com and www.googlecom
func MissingDot(domain, tld string) []string {
	var res []string
	nbDot := strings.Count(domain, ".")
	for i := 0; i < nbDot; i++ {
		res = append(res, concatDomainTopLevel(strings.Replace(domain, ".", "", 1), tld))
	}
	return res
}

// Homoglyphs One or more characters that look similar to another character but are different are called homogylphs. An example is that the lower case l looks similar to the numeral one, e.g. l vs 1. For example, google.com becomes goog1e.com.
func ReplaceByHomoglyphs(domain, tld string) []string {
	var res []string
	domainTmp := domain
	for _, c := range domainTmp {
		if homoglyphs, ok := homoglyphMap[c]; ok {
			for _, homoglyph := range homoglyphs {
				res = append(res, concatDomainTopLevel(strings.Replace(domain, string(c), string(homoglyph), 1), tld))
			}
		}
	}
	return res
}

// change top domain -- data.iana.org/TLD/tlds-alpha-by-domain.txt
// For example, www.trademe.co.nz becomes www.trademe.co.nz and www.google.com becomes www.google.org
func ChangeTopDomain(domain string) []string {
	var res []string
	for _, tld := range topDomain {
		res = append(res, concatDomainTopLevel(domain, tld))
	}
	return res
}

func concatDomainTopLevel(domain, tld string) string {
	return domain + "." + tld
}
