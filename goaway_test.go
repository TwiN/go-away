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

func TestSneakyBadWords(t *testing.T)  {
	words := []string{"A$$", "4ss", "4s$", "a S s", "a $ s"}

	for _, w := range words {
		if !goaway.IsProfane(w) {
			t.Error("Expected true, got false from word", w)
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