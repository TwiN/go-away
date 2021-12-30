package goaway

// DefaultCharacterReplacements is the mapping of all characters that are replaced by other characters before
// attempting to find a profanity.
var DefaultCharacterReplacements = map[rune]rune{
	// Leetspeak
	'0': 'o',
	'1': 'i',
	'3': 'e',
	'4': 'a',
	'5': 's',
	'7': 'l',
	'$': 's',
	'!': 'i',
	'+': 't',
	'#': 'h',
	'@': 'a',
	'<': 'c',
	// Special characters
	'-': ' ',
	'_': ' ',
	'|': ' ',
	'.': ' ',
	',': ' ',
	'(': ' ',
	')': ' ',
	'>': ' ',
	'"': ' ',
	'`': ' ',
	'~': ' ',
	'*': ' ',
	'&': ' ',
	'%': ' ',
	'?': ' ',
}
