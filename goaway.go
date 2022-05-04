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
	sanitizeSpecialCharacters bool // Whether to replace characters with the value ' ' in characterReplacements
	sanitizeLeetSpeak         bool // Whether to replace characters with a non-' ' value in characterReplacements
	sanitizeAccents           bool
	sanitizeSpaces            bool

	profanities    []string
	falseNegatives []string
	falsePositives []string

	characterReplacements map[rune]rune
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
		characterReplacements:     DefaultCharacterReplacements,
	}
}

// WithSanitizeLeetSpeak allows configuring whether the sanitization process should also take into account leetspeak
//
// Leetspeak characters are characters to be replaced by non-' ' values in the characterReplacements map.
// For instance, '4' is replaced by 'a' and '3' is replaced by 'e', which means that "4sshol3" would be
// sanitized to "asshole", which would be detected as a profanity.
//
// By default, this is set to true.
func (g *ProfanityDetector) WithSanitizeLeetSpeak(sanitize bool) *ProfanityDetector {
	g.sanitizeLeetSpeak = sanitize
	return g.buildCharacterReplacements()
}

// WithSanitizeSpecialCharacters allows configuring whether the sanitization process should also take into account
// special characters.
//
// Special characters are characters that are part of the characterReplacements map (DefaultCharacterReplacements by
// default) and are to be removed during the sanitization step.
//
// For instance, "fu_ck" would be sanitized to "fuck", which would be detected as a profanity.
//
// By default, this is set to true.
func (g *ProfanityDetector) WithSanitizeSpecialCharacters(sanitize bool) *ProfanityDetector {
	g.sanitizeSpecialCharacters = sanitize
	return g.buildCharacterReplacements()
}

// WithSanitizeAccents allows configuring of whether the sanitization process should also take into account accents.
// By default, this is set to true, but since this adds a bit of overhead, you may disable it if your use case
// is time-sensitive or if the input doesn't involve accents (i.e. if the input can never contain special characters)
func (g *ProfanityDetector) WithSanitizeAccents(sanitize bool) *ProfanityDetector {
	g.sanitizeAccents = sanitize
	return g
}

// WithSanitizeSpaces allows configuring whether the sanitization process should also take into account spaces
func (g *ProfanityDetector) WithSanitizeSpaces(sanitize bool) *ProfanityDetector {
	g.sanitizeSpaces = sanitize
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

// WithCustomCharacterReplacements allows configuring characters that to be replaced by other characters.
//
// Note that all entries that have the value ' ' are considered as special characters while all entries with a value
// that is not ' ' are considered as leet speak.
//
// Defaults to DefaultCharacterReplacements
func (g *ProfanityDetector) WithCustomCharacterReplacements(characterReplacements map[rune]rune) *ProfanityDetector {
	g.characterReplacements = characterReplacements
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

func (g *ProfanityDetector) indexToRune(s string, index int) int {
	count := 0
	for i := range s {
		if i == index {
			return count
		}
		count++
	}
	return -1
}

// Censor takes in a string (word or sentence) and tries to censor all profanities found.
func (g *ProfanityDetector) Censor(s string) string {
	censored := []rune(s)
	var originalIndexes []int
	s, originalIndexes = g.sanitize(s, true)
	// Check for false negatives
	for _, word := range g.falseNegatives {
		currentIndex := 0
		for currentIndex != -1 {
			if foundIndex := strings.Index(s[currentIndex:], word); foundIndex != -1 {
				for i := 0; i < len([]rune(word)); i++ {
					runeIndex := g.indexToRune(string(censored), currentIndex+foundIndex+i)
					censored[originalIndexes[runeIndex]] = '*'
				}
				currentIndex += foundIndex + len([]rune(word))
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
				foundRuneIndex := g.indexToRune(s, foundIndex)
				originalIndexes = append(originalIndexes[:foundRuneIndex], originalIndexes[foundRuneIndex+len(word):]...)
				currentIndex += foundIndex + len([]rune(word))
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
				for i := 0; i < len([]rune(word)); i++ {
					runeIndex := g.indexToRune(string(censored), currentIndex+foundIndex+i)
					censored[originalIndexes[runeIndex]] = '*'
				}
				currentIndex += foundIndex + len([]rune(word))
			} else {
				break
			}
		}
	}
	return string(censored)
}

func (g ProfanityDetector) sanitize(s string, rememberOriginalIndexes bool) (string, []int) {
	s = strings.ToLower(s)
	if g.sanitizeLeetSpeak && !rememberOriginalIndexes && g.sanitizeSpecialCharacters {
		s = strings.ReplaceAll(s, "()", "o")
	}
	sb := strings.Builder{}
	for _, char := range s {
		if replacement, found := g.characterReplacements[char]; found {
			if g.sanitizeSpecialCharacters && replacement == ' ' {
				// If the replacement is a space, and we're sanitizing special characters speak, we replace.
				sb.WriteRune(replacement)
				continue
			} else if g.sanitizeLeetSpeak && replacement != ' ' {
				// If the replacement isn't a space, and we're sanitizing leet speak, we replace.
				sb.WriteRune(replacement)
				continue
			}
		}
		sb.WriteRune(char)
	}
	s = sb.String()
	if g.sanitizeAccents {
		s = removeAccents(s)
	}
	var originalIndexes []int
	if rememberOriginalIndexes {
		for i, c := range []rune(s) {
			// If spaces aren't being sanitized, appending to the original indices prevents off-by-one errors later on.
			if c != ' ' || !g.sanitizeSpaces {
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

// buildCharacterReplacements builds characterReplacements if WithSanitizeLeetSpeak or WithSanitizeSpecialCharacters is
// called.
//
// If this is not called, DefaultCharacterReplacements
func (g *ProfanityDetector) buildCharacterReplacements() *ProfanityDetector {
	g.characterReplacements = make(map[rune]rune)
	if g.sanitizeSpecialCharacters {
		g.characterReplacements['-'] = ' '
		g.characterReplacements['_'] = ' '
		g.characterReplacements['|'] = ' '
		g.characterReplacements['.'] = ' '
		g.characterReplacements[','] = ' '
		g.characterReplacements['('] = ' '
		g.characterReplacements[')'] = ' '
		g.characterReplacements['<'] = ' '
		g.characterReplacements['>'] = ' '
		g.characterReplacements['"'] = ' '
		g.characterReplacements['`'] = ' '
		g.characterReplacements['~'] = ' '
		g.characterReplacements['*'] = ' '
		g.characterReplacements['&'] = ' '
		g.characterReplacements['%'] = ' '
		g.characterReplacements['$'] = ' '
		g.characterReplacements['#'] = ' '
		g.characterReplacements['@'] = ' '
		g.characterReplacements['!'] = ' '
		g.characterReplacements['?'] = ' '
		g.characterReplacements['+'] = ' '
	}
	if g.sanitizeLeetSpeak {
		g.characterReplacements['4'] = 'a'
		g.characterReplacements['$'] = 's'
		g.characterReplacements['!'] = 'i'
		g.characterReplacements['+'] = 't'
		g.characterReplacements['#'] = 'h'
		g.characterReplacements['@'] = 'a'
		g.characterReplacements['0'] = 'o'
		g.characterReplacements['1'] = 'i'
		g.characterReplacements['7'] = 'l'
		g.characterReplacements['3'] = 'e'
		g.characterReplacements['5'] = 's'
		g.characterReplacements['<'] = 'c'
	}
	return g
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
