package goaway

// FalsePositives is a list of words that may wrongly trigger the profanity filter
var FalsePositives = []string{}

// BlockedWords is a set of blocked words derived from Profanities and FalsePositives
var BlockedWords = make(map[string]struct{})

func init() {
	for _, word := range Profanities {
		BlockedWords[word] = struct{}{}
	}

	for _, word := range FalsePositives {
		BlockedWords[word] = struct{}{}
	}
}

// IsBlockedWord checks if a word is blocked
func IsBlockedWord(word string) bool {
	_, blocked := BlockedWords[word]
	return blocked
}

// FilterBlockedWords filters out blocked words from a given list of words
func FilterBlockedWords(words []string) []string {
	filtered := make([]string, 0, len(words))
	for _, word := range words {
		if !IsBlockedWord(word) {
			filtered = append(filtered, word)
		}
	}
	return filtered
}
