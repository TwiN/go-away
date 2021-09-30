package goaway

import (
	"strings"
	"unicode"

	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

const (
	space              = " "
	firstRuneSupported = ' '
	lastRuneSupported  = '~'
)

var (
	defaultProfanityDetector *ProfanityDetector
	removeAccentsTransformer transform.Transformer
)

// ProfanityDetector contains the dictionaries as well as the configuration
// for determining how profanity detection is handled
type ProfanityDetector struct {
	sanitizeSpecialCharacters bool
	sanitizeLeetSpeak         bool
	sanitizeAccents           bool
	sanitizeSpaces            bool

	profanities    []string
	falseNegatives []string
	falsePositives []string
}

// NewProfanityDetector creates a new ProfanityDetector
func NewProfanityDetector() *ProfanityDetector {
	return &ProfanityDetector{
		sanitizeSpecialCharacters: true,
		sanitizeLeetSpeak:         true,
		sanitizeAccents:           true,
		sanitizeSpaces:            true,
		profanities:               DefaultProfanities,
		falsePositives:            DefaultFalsePositives,
		falseNegatives:            DefaultFalseNegatives,
	}
}

// WithSanitizeLeetSpeak allows configuring whether the sanitization process should also take into account leetspeak
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

// WithCustomDictionary allows configuring whether the sanitization process should also take into account
// custom profanities, false positives and false negatives dictionaries.
func (g *ProfanityDetector) WithCustomDictionary(profanities, falsePositives, falseNegatives []string) *ProfanityDetector {
	g.profanities = profanities
	g.falsePositives = falsePositives
	g.falseNegatives = falseNegatives
	return g
}

// WithSanitizeSpaces allows configuring whether the sanitization process should also take into account spaces
func (g *ProfanityDetector) WithSanitizeSpaces(sanitize bool) *ProfanityDetector {
	g.sanitizeSpaces = sanitize
	return g
}

// IsProfane takes in a string (word or sentence) and look for profanities.
// Returns a boolean
func (g *ProfanityDetector) IsProfane(s string) bool {
	return len(g.ExtractProfanity(s)) > 0
}

// ExtractProfanity takes in a string (word or sentence) and look for profanities.
// Returns the first profanity found, or an empty string if none are found
func (g *ProfanityDetector) ExtractProfanity(s string) string {
	s = g.sanitize(s)
	// Check for false negatives
	for _, word := range g.falseNegatives {
		if match := strings.Contains(s, word); match {
			return word
		}
	}
	// Remove false positives
	for _, word := range g.falsePositives {
		s = strings.Replace(s, word, "", -1)
	}
	// Check for profanities
	for _, word := range g.profanities {
		if match := strings.Contains(s, word); match {
			return word
		}
	}
	return ""
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
	if g.sanitizeSpaces {
		s = strings.Replace(s, space, "", -1)
	}
	if g.sanitizeAccents {
		s = removeAccents(s)
	}
	return s
}

// removeAccents strips all accents from characters.
// Only called if ProfanityDetector.removeAccents is set to true
func removeAccents(s string) string {
	if removeAccentsTransformer == nil {
		removeAccentsTransformer = transform.Chain(norm.NFD, runes.Remove(runes.In(unicode.Mn)), norm.NFC)
	}
	for _, character := range s {
		// If there's a character outside the range of supported runes, there might be some accented words
		if character < firstRuneSupported || character > lastRuneSupported {
			s, _, _ = transform.String(removeAccentsTransformer, s)
			break
		}
	}
	return s
}

// IsProfane checks whether there are any profanities in a given string (word or sentence).
//
// Uses the default ProfanityDetector
func IsProfane(s string) bool {
	if defaultProfanityDetector == nil {
		defaultProfanityDetector = NewProfanityDetector()
	}
	return defaultProfanityDetector.IsProfane(s)
}

// ExtractProfanity takes in a string (word or sentence) and look for profanities.
// Returns the first profanity found, or an empty string if none are found
//
// Uses the default ProfanityDetector
func ExtractProfanity(s string) string {
	if defaultProfanityDetector == nil {
		defaultProfanityDetector = NewProfanityDetector()
	}
	return defaultProfanityDetector.ExtractProfanity(s)
}
