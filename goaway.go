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

	profanities                     []string
	falseNegatives                  []string
	falsePositives                  []string
	specialCharactersReplacementMap map[rune]rune
	leetSpeekReplacementMap         map[rune]rune
}

// NewProfanityDetector creates a new ProfanityDetector
func NewProfanityDetector() *ProfanityDetector {
	return &ProfanityDetector{
		sanitizeSpecialCharacters:       true,
		sanitizeLeetSpeak:               true,
		sanitizeAccents:                 true,
		sanitizeSpaces:                  true,
		profanities:                     DefaultProfanities,
		falsePositives:                  DefaultFalsePositives,
		falseNegatives:                  DefaultFalseNegatives,
		specialCharactersReplacementMap: createIgnoreMap(DefaultIgnoredSpecialCharacters),
		leetSpeekReplacementMap:         DefaultLeetspeekCharactersReplacement,
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

// WithSpecialCharacters allows configuring special characters that should be removed before checking for profanities
func (g *ProfanityDetector) WithSpecialCharacters(specialCharacters []rune) *ProfanityDetector {
	g.specialCharactersReplacementMap = createIgnoreMap(specialCharacters)
	return g
}

// WithLeetSpeakReplacements allows configuring custom leet speak replacements
func (g *ProfanityDetector) WithLeetSpeakReplacements(replacementMap map[rune]rune) *ProfanityDetector {
	g.leetSpeekReplacementMap = replacementMap
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
	s, _ = g.sanitize(s, false)
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

// Censor takes in a string (word or sentence) and tries to censor all profanities found.
func (g *ProfanityDetector) Censor(s string) string {
	censored := s
	var originalIndexes []int
	s, originalIndexes = g.sanitize(s, true)
	// Check for false negatives
	for _, word := range g.falseNegatives {
		currentIndex := 0
		for currentIndex != -1 {
			if foundIndex := strings.Index(s[currentIndex:], word); foundIndex != -1 {
				for i := 0; i < len(word); i++ {
					censored = censored[:originalIndexes[foundIndex+currentIndex+i]] + "*" + censored[originalIndexes[foundIndex+currentIndex+i]+1:]
				}
				currentIndex += foundIndex + len(word)
			} else {
				break
			}
		}
	}
	// Remove false positives
	for _, word := range g.falsePositives {
		currentIndex := 0
		for currentIndex != -1 {
			if foundIndex := strings.Index(s[currentIndex:], word); foundIndex != -1 {
				originalIndexes = append(originalIndexes[:foundIndex+currentIndex], originalIndexes[foundIndex+len(word):]...)
				currentIndex += foundIndex + len(word)
			} else {
				break
			}
		}
		s = strings.Replace(s, word, "", -1)
	}
	// Check for profanities
	for _, word := range g.profanities {
		currentIndex := 0
		for currentIndex != -1 {
			if foundIndex := strings.Index(s[currentIndex:], word); foundIndex != -1 {
				for i := 0; i < len(word); i++ {
					censored = censored[:originalIndexes[foundIndex+currentIndex+i]] + "*" + censored[originalIndexes[foundIndex+currentIndex+i]+1:]
				}
				currentIndex += foundIndex + len(word)
			} else {
				break
			}
		}
	}
	return censored
}

func (g ProfanityDetector) sanitize(s string, rememberOriginalIndexes bool) (string, []int) {
	s = strings.ToLower(s)
	if g.sanitizeLeetSpeak && !rememberOriginalIndexes && g.sanitizeSpecialCharacters {
		s = strings.ReplaceAll(s, "()", "o")
	}
	sb := strings.Builder{}
	for _, char := range s {
		replaced := false
		if g.sanitizeLeetSpeak {

			_, isSpecialCharacter := g.specialCharactersReplacementMap[char]
			if !isSpecialCharacter || g.sanitizeSpecialCharacters {
				repl, found := g.leetSpeekReplacementMap[char]
				if found {
					sb.WriteRune(repl)
					replaced = true
				}
			}
		}
		if g.sanitizeSpecialCharacters {
			if !replaced {
				repl, found := g.specialCharactersReplacementMap[char]
				if found {
					sb.WriteRune(repl)
					replaced = true
				}
			}
		}
		if !replaced {
			sb.WriteRune(char)
		}
	}
	s = sb.String()

	if g.sanitizeAccents {
		s = removeAccents(s)
	}
	var originalIndexes []int
	if rememberOriginalIndexes {
		for i, c := range s {
			if c != ' ' {
				originalIndexes = append(originalIndexes, i)
			}
		}
	}
	if g.sanitizeSpaces {
		s = strings.Replace(s, space, "", -1)
	}
	return s, originalIndexes
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

// Censor takes in a string (word or sentence) and tries to censor all profanities found.
//
// Uses the default ProfanityDetector
func Censor(s string) string {
	if defaultProfanityDetector == nil {
		defaultProfanityDetector = NewProfanityDetector()
	}
	return defaultProfanityDetector.Censor(s)
}
