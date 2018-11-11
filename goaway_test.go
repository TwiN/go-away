package goaway_test

import (
	"github.com/TwinProduction/go-away"
	"testing"
)

func TestBadWords(t *testing.T)  {
	words := []string{"fuck", "ass", "poop", "penis", "bitch"}

	for _, w := range words {
		if !goaway.IsProfane(w) {
			t.Error("Expected true, got false from word", w)
		}
	}
}

func TestSentencesWithBadWords(t *testing.T)  {
	sentences := []string{"What the fuck is your problem", "Go away, asshole!"}

	for _, s := range sentences {
		if !goaway.IsProfane(s) {
			t.Error("Expected true, got false from sentence", s)
		}
	}
}

func TestSneakyBadWords(t *testing.T)  {
	words := []string{"A$$", "4ss", "4s$", "a S s", "a $ s", "@$$h073", "f    u     c k"}

	for _, w := range words {
		if !goaway.IsProfane(w) {
			t.Error("Expected true, got false from word", w)
		}
	}
}

func TestSentencesWithSneakyBadWords(t *testing.T)  {
	sentences := []string{"You smell p00p", "Go away, a$$h0l3!"}

	for _, s := range sentences {
		if !goaway.IsProfane(s) {
			t.Error("Expected true, got false from sentence", s)
		}
	}
}

func TestNormalWords(t *testing.T)  {
	words := []string{"hello", "world", "whats", "up"}

	for _, w := range words {
		if goaway.IsProfane(w) {
			t.Error("Expected false, got false from word", w)
		}
	}
}

func TestSentencesWithNoProfanities(t *testing.T)  {
	sentences := []string{"hello, my friend", "what's up?", "do you want to play bingo?"}

	for _, s := range sentences {
		if goaway.IsProfane(s) {
			t.Error("Expected false, got false from sentence", s)
		}
	}
}