package goaway

import (
	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
	"strings"
	"unicode"
)

const Space = " "

// Takes in a string (word or sentence) and look for profanities.
// Returns a boolean
func IsProfane(s string) bool {
	s = strings.Replace(sanitize(s), Space, "", -1) // Sanitize leetspeak AND remove all spaces
	for _, falsePositive := range falsePositives {
		s = strings.Replace(s, falsePositive, "", -1)
	}
	for _, word := range profanities {
		if match := strings.Contains(s, word); match {
			return true
		}
	}
	return false
}

func sanitize(s string) string {
	s = removeAccents(s)
	s = strings.ToLower(s)
	s = strings.Replace(s, "0", "o", -1)
	s = strings.Replace(s, "1", "i", -1)
	s = strings.Replace(s, "3", "e", -1)
	s = strings.Replace(s, "4", "a", -1)
	s = strings.Replace(s, "5", "s", -1)
	s = strings.Replace(s, "6", "b", -1)
	s = strings.Replace(s, "7", "l", -1)
	s = strings.Replace(s, "8", "b", -1)
	s = strings.Replace(s, "@", "a", -1)
	s = strings.Replace(s, "!", "i", -1)
	s = strings.Replace(s, "+", "t", -1)
	s = strings.Replace(s, "$", "s", -1)
	s = strings.Replace(s, "()", "o", -1)
	s = strings.Replace(s, "_", Space, -1)
	s = strings.Replace(s, "-", Space, -1)
	s = strings.Replace(s, "*", Space, -1)
	return s
}

func removeAccents(s string) string {
	t := transform.Chain(norm.NFD, runes.Remove(runes.In(unicode.Mn)), norm.NFC)
	output, _, e := transform.String(t, s)
	if e != nil {
		panic(e)
	}
	return output
}
