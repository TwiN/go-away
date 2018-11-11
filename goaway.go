package goaway

import (
	"log"
	"io/ioutil"
	"regexp"
	"strings"
)

var profanities []string
var initialized bool

func Initialize(finished chan bool) {
	if initialized {
		finished <- true
		return // already initialized
	}
	log.Println("[Initialize] Initializing go-away")
	b, err := ioutil.ReadFile("profanities.txt")
	if err != nil {
		log.Fatalln("[Initialize] Error reading profanities file:", err.Error())
	}

	for _, profanity := range strings.Split(string(b), "\n") {
		profanities = append(profanities, profanity)
	}
	initialized = true
	finished <- true
}

/**
 * Takes in a string (word or sentence) and look for profanities.
 * Returns a boolean
 */
func IsProfane(s string) bool {
	if !initialized {
		log.Println("You must call goaway.Initialize() first")
		finished := make(chan bool)
		Initialize(finished)
		<- finished
	}
	s = strings.Replace(sanitize(s), " ", "", -1) // Sanitize leetspeak AND remove all spaces
	for _, word := range profanities {
		wordPattern := `\b` + word + `\b`
		match, _ := regexp.MatchString(wordPattern, s)

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
