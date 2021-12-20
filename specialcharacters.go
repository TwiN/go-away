package goaway

var DefaultSpecialCharacters = []rune("._-?()|~")

func createReplacementMap(runes []rune) map[rune]rune {
	specialReplacementMap := make(map[rune]rune, len(runes))
	for _, character := range runes {
		specialReplacementMap[character] = ' '
	}
	return specialReplacementMap
}

// DefaultLeetspeekCharactersReplacement contains list of special character runes that will be turned to special character mapping
var DefaultSpecialCharacterReplacements = createReplacementMap(DefaultSpecialCharacters)

// DefaultLeetspeekCharactersReplacement contains mapping for leet speak characters mapping. Note that special characters will not be mapped if sanitizeSpecialCharacters is set to false
var DefaultLeetspeekCharactersReplacement = map[rune]rune{
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
