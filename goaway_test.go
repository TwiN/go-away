package goaway_test

import (
	"github.com/mattwhite180/go-away/v1"
	"testing"
)

func TestBadWords(t *testing.T) {
	words := []string{"fuck", "ass", "shit", "penis", "bitch"}
	for _, w := range words {
		if !goaway.IsProfane(w) {
			t.Error("Expected true, got false from word '", w, "'")
		}
	}
}

func TestBadWordsWithAccentedLetters(t *testing.T) {
	words := []string{"fučk", "ÄšŚ", "sh1t", "pÉnìŚ", "bitčh"}
	for _, w := range words {
		if !goaway.IsProfane(w) {
			t.Error("Expected true, got false from word '", w, "'")
		}
	}
}

func TestSentencesWithBadWords(t *testing.T) {
	sentences := []string{"What the fuck is your problem", "Go away, asshole!"}
	for _, s := range sentences {
		if !goaway.IsProfane(s) {
			t.Error("Expected true, got false from sentence '", s, "'")
		}
	}
}

func TestSneakyBadWords(t *testing.T) {
	words := []string{"A$$", "4ss", "4s$", "a S s", "a $ s", "@$$h073", "f    u     c k"}
	for _, w := range words {
		if !goaway.IsProfane(w) {
			t.Error("Expected true, got false from word '", w, "'")
		}
	}
}

func TestSentencesWithSneakyBadWords(t *testing.T) {
	sentences := []string{"You smell $h1t", "Go away, a$$h0l3!"}
	for _, s := range sentences {
		if !goaway.IsProfane(s) {
			t.Error("Expected true, got false from sentence '", s, "'")
		}
	}
}

func TestNormalWords(t *testing.T) {
	words := []string{"hello", "world", "whats", "up"}
	for _, w := range words {
		if goaway.IsProfane(w) {
			t.Error("Expected false, got true from word '", w, "'")
		}
	}
}

func TestSentencesWithNoProfanities(t *testing.T) {
	sentences := []string{"hello, my friend", "what's up?", "do you want to play bingo?"}
	for _, s := range sentences {
		if goaway.IsProfane(s) {
			t.Error("Expected false, got true from sentence '", s, "'")
		}
	}
}

func TestSentencesWithFalseProfanities(t *testing.T) {
	sentences := []string{"I am from Scuntthorpe, north Lincolnshire", "He is an associate of mine"}
	for _, s := range sentences {
		if goaway.IsProfane(s) {
			t.Error("Expected false, got true from sentence '", s, "'")
		}
	}
}

func TestSentencesWithFalsePositivesAndProfanities(t *testing.T) {
	sentences := []string{"You are a poopy associate", "Go back to Scuntthorpe, Asshole!"}
	for _, s := range sentences {
		if !goaway.IsProfane(s) {
			t.Error("Expected true, got false from sentence '", s, "'")
		}
	}
}
