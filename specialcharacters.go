package goaway

var DefaultIgnoredCharacters = []rune("._-?()|~")

func createIgnoreMap(runes []rune) map[rune]rune {
	specialReplacementMap := make(map[rune]rune, len(runes))
	for _, character := range runes {
		specialReplacementMap[character] = ' '
	}
	return specialReplacementMap
}

// DefaultLeetSpeakCharactersReplacement contains mapping for leet speak characters.
// Note that special leet speak characters will not be mapped if sanitizeSpecialCharacters is set to false
var DefaultLeetSpeakCharactersReplacement = map[rune]rune{
	'4': 'a',
	'$': 's',
	'!': 'i',
	'+': 't',
	'#': 'h',
	'@': 'a',
	'0': 'o',
	'1': 'i',
	'7': 'l',
	'3': 'e',
	'5': 's',
}
