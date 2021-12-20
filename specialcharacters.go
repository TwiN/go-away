package goaway

var DefaultSpecialCharacters = []rune("._-?()|~")

func createReplacementMap(runes []rune) map[rune]rune {
	specialReplacementMap := make(map[rune]rune, len(runes))
	for _, character := range runes {
		specialReplacementMap[character] = ' '
	}
	return specialReplacementMap
}

var DefaultSpecialCharacterReplacements = createReplacementMap(DefaultSpecialCharacters)

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
