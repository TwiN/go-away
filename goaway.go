package goaway

import (
	"strings"
)

var profanities = []string{"anal","anus","arse","ass","ballsack","balls","bastard","bitch","biatch","bloody","blowjob","bollock","bollok","boner","boob","bugger","bum","butt","clitoris","cock","coon","crap","cunt","dick","dildo","dyke","fag","feck","fellate","fellatio","felching","fuck","fudgepacker","flange","homo","jerk","jizz","labia","muff","naked","nigger","nigga","nude","penis","piss","poop","porn","prick","pube","pussy","queer","scrotum","sex","shit","slut","spunk","suckmy","tit","tosser","turd","twat","vagina","wank","whore"}

/**
 * Takes in a string (word or sentence) and look for profanities.
 * Returns a boolean
 */
func IsProfane(s string) bool {
	s = strings.Replace(sanitize(s), " ", "", -1) // Sanitize leetspeak AND remove all spaces
	for _, word := range profanities {
		match := strings.Contains(s, word)
		if match {
			return true
		}
	}
	return false
}

func sanitize(s string) string {
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
	s = strings.Replace(s, "!", "a", -1)
	s = strings.Replace(s, "$", "s", -1)
	s = strings.Replace(s, "()", "o", -1)
	s = strings.Replace(s, "_", " ", -1)
	s = strings.Replace(s, "-", " ", -1)
	s = strings.Replace(s, "*", " ", -1)
	return s
}
