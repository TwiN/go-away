package goaway

import (
	"strings"
	"unicode"

	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

const Space = " "

var (
	defaultProfanityDetector *ProfanityDetector
	removeAccentsTransformer transform.Transformer
)

// ProfanityDetector
type ProfanityDetector struct {
	sanitizeSpecialCharacters bool
	sanitizeLeetSpeak         bool
	sanitizeAccents           bool
}

// NewProfanityDetector creates a new ProfanityDetector
func NewProfanityDetector() *ProfanityDetector {
	return &ProfanityDetector{
		sanitizeSpecialCharacters: true,
		sanitizeLeetSpeak:         true,
		sanitizeAccents:           true,
	}
}

// WithSanitizeLeetSpeak allows configuring whether the sanitization process should also take into account
// leetspeak
func (g *ProfanityDetector) WithSanitizeLeetSpeak(sanitize bool) *ProfanityDetector {
	g.sanitizeLeetSpeak = sanitize
	return g
}

// WithSanitizeSpecialCharacters allows configuring whether the sanitization process should also take into account
// special characters
func (g *ProfanityDetector) WithSanitizeSpecialCharacters(sanitize bool) *ProfanityDetector {
	g.sanitizeSpecialCharacters = sanitize
	return g
}

// WithSanitizeAccents allows configuring of whether the sanitization process should also take into account accents.
// By default, this is set to true, but since this adds a bit of overhead, you may disable it if your use case
// is time-sensitive or if the input doesn't involve accents (i.e. if the input can never contain special characters)
func (g *ProfanityDetector) WithSanitizeAccents(sanitize bool) *ProfanityDetector {
	g.sanitizeAccents = sanitize
	return g
}

// IsProfane takes in a string (word or sentence) and look for profanities.
// Returns a boolean
func (g *ProfanityDetector) IsProfane(s string) bool {
	s = g.sanitize(s)
	// Check for false false positives
	for _, word := range falseNegatives {
		if match := strings.Contains(s, word); match {
			return true
		}
	}
	// Remove false positives
	for _, word := range falsePositives {
		s = strings.Replace(s, word, "", -1)
	}
	// Check for profanities
	for _, word := range profanities {
		if match := strings.Contains(s, word); match {
			return true
		}
	}
	return false
}

func (g ProfanityDetector) sanitize(s string) string {
	s = strings.ToLower(s)
	if g.sanitizeLeetSpeak {
		s = strings.Replace(s, "0", "o", -1)
		s = strings.Replace(s, "1", "i", -1)
		s = strings.Replace(s, "3", "e", -1)
		s = strings.Replace(s, "4", "a", -1)
		s = strings.Replace(s, "5", "s", -1)
		s = strings.Replace(s, "6", "b", -1)
		s = strings.Replace(s, "7", "l", -1)
		s = strings.Replace(s, "8", "b", -1)
	}
	if g.sanitizeSpecialCharacters {
		if g.sanitizeLeetSpeak {
			s = strings.Replace(s, "@", "a", -1)
			s = strings.Replace(s, "+", "t", -1)
			s = strings.Replace(s, "$", "s", -1)
			s = strings.Replace(s, "#", "h", -1)
			s = strings.Replace(s, "()", "o", -1)
		}
		s = strings.Replace(s, "_", "", -1)
		s = strings.Replace(s, "-", "", -1)
		s = strings.Replace(s, "*", "", -1)
		s = strings.Replace(s, "'", "", -1)
		s = strings.Replace(s, "?", "", -1)
		s = strings.Replace(s, "!", "", -1)
	}
	s = strings.Replace(s, Space, "", -1)
	if g.sanitizeAccents {
		s = removeAccents(s)
	}
	return s
}

func removeAccents(s string) string {
	if removeAccentsTransformer == nil {
		removeAccentsTransformer = transform.Chain(norm.NFD, runes.Remove(runes.In(unicode.Mn)), norm.NFC)
	}
	output, _, _ := transform.String(removeAccentsTransformer, s)
	return output
}

// IsProfane checks whether there are any profanities in a given string (word or sentence).
// Uses the default ProfanityDetector
func IsProfane(s string) bool {
	if defaultProfanityDetector == nil {
		defaultProfanityDetector = NewProfanityDetector()
	}
	return defaultProfanityDetector.IsProfane(s)
}
