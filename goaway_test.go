package goaway_test

import (
	"github.com/TwinProduction/go-away"
	"testing"
)

func TestBadWords(t *testing.T) {
	words := []string{"fuck", "ass", "poop", "penis", "bitch"}
	for _, w := range words {
		if !goaway.IsProfane(w) {
			t.Error("Expected true, got false from word", w)
		}
	}
}

func TestBadWordsWithAccentedLetters(t *testing.T) {
	words := []string{"fučk", "ÄšŚ", "pÓöp", "pÉnìŚ", "bitčh"}
	for _, w := range words {
		if !goaway.IsProfane(w) {
			t.Error("Expected true, got false from word", w)
		}
	}
}

func TestSentencesWithBadWords(t *testing.T) {
	sentences := []string{"What the fuck is your problem", "Go away, asshole!"}
	for _, s := range sentences {
		if !goaway.IsProfane(s) {
			t.Error("Expected true, got false from sentence", s)
		}
	}
}

func TestSneakyBadWords(t *testing.T) {
	words := []string{"A$$", "4ss", "4s$", "a S s", "a $ s", "@$$h073", "f    u     c k", "4r5e", "5h1t", "5hit", "a55", "ar5e", "a_s_s", "b!tch", "b!+ch"}
	for _, w := range words {
		if !goaway.IsProfane(w) {
			t.Error("Expected true, got false from word", w)
		}
	}
}

func TestSentencesWithSneakyBadWords(t *testing.T) {
	sentences := []string{"You smell p00p", "Go away, a$$h0l3!"}
	for _, s := range sentences {
		if !goaway.IsProfane(s) {
			t.Error("Expected true, got false from sentence", s)
		}
	}
}

func TestNormalWords(t *testing.T) {
	words := []string{"hello", "world", "whats", "up"}
	for _, w := range words {
		if goaway.IsProfane(w) {
			t.Error("Expected false, got true from word", w)
		}
	}
}

func TestSentencesWithNoProfanities(t *testing.T) {
	sentences := []string{"hello, my friend", "what's up?", "do you want to play bingo?"}
	for _, s := range sentences {
		if goaway.IsProfane(s) {
			t.Error("Expected false, got false from sentence", s)
		}
	}
}

func TestSentencesWithFalsePositives(t *testing.T) {
	sentences := []string{"I am from Scuntthorpe, north Lincolnshire", "He is an associate of mine", "Are you an assassin?"}
	for _, s := range sentences {
		if goaway.IsProfane(s) {
			t.Error("Expected false, got true from sentence", s)
		}
	}
}

func TestSentencesWithFalsePositivesAndProfanities(t *testing.T) {
	sentences := []string{"You are a shitty associate", "Go away, asshole!"}
	for _, s := range sentences {
		if !goaway.IsProfane(s) {
			t.Error("Expected true, got false from sentence", s)
		}
	}
}
