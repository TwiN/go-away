package goaway

import (
	"testing"
)

func BenchmarkIsProfaneWhenShortStringHasNoProfanity(b *testing.B) {
	for n := 0; n < b.N; n++ {
		IsProfane("aaaaaaaaaaaaaa")
	}
}

func BenchmarkIsProfaneWhenShortStringHasProfanityAtTheStart(b *testing.B) {
	for n := 0; n < b.N; n++ {
		IsProfane("fuckaaaaaaaaaa")
	}
}

func BenchmarkIsProfaneWhenShortStringHasProfanityInTheMiddle(b *testing.B) {
	for n := 0; n < b.N; n++ {
		IsProfane("aaaaafuckaaaaa")
	}
}

func BenchmarkIsProfaneWhenShortStringHasProfanityAtTheEnd(b *testing.B) {
	for n := 0; n < b.N; n++ {
		IsProfane("aaaaaaaaaafuck")
	}
}

func BenchmarkIsProfaneWhenMediumStringHasNoProfanity(b *testing.B) {
	for n := 0; n < b.N; n++ {
		IsProfane("How are you doing today?")
	}
}

func BenchmarkIsProfaneWhenMediumStringHasProfanityAtTheStart(b *testing.B) {
	for n := 0; n < b.N; n++ {
		IsProfane("Shit, you're cute today.")
	}
}

func BenchmarkIsProfaneWhenMediumStringHasProfanityInTheMiddle(b *testing.B) {
	for n := 0; n < b.N; n++ {
		IsProfane("How are you fu ck doing?")
	}
}

func BenchmarkIsProfaneWhenMediumStringHasProfanityAtTheEnd(b *testing.B) {
	for n := 0; n < b.N; n++ {
		IsProfane("you're cute today. Fuck.")
	}
}

func BenchmarkIsProfaneWhenLongStringHasNoProfanity(b *testing.B) {
	for n := 0; n < b.N; n++ {
		IsProfane("Hello John Doe, I hope you're feeling well, as I come today bearing terrible news regarding your favorite chocolate chip cookie brand")
	}
}

func BenchmarkIsProfaneWhenLongStringHasProfanityAtTheStart(b *testing.B) {
	for n := 0; n < b.N; n++ {
		IsProfane("Fuck John Doe, I hope you're feeling well, as I come today bearing terrible news regarding your favorite chocolate chip cookie brand")
	}
}

func BenchmarkIsProfaneWhenLongStringHasProfanityInTheMiddle(b *testing.B) {
	for n := 0; n < b.N; n++ {
		IsProfane("Hello John Doe, I hope you're feeling well, as I come today bearing shitty news regarding your favorite chocolate chip cookie brand")
	}
}

func BenchmarkIsProfaneWhenLongStringHasProfanityAtTheEnd(b *testing.B) {
	for n := 0; n < b.N; n++ {
		IsProfane("Hello John Doe, I hope you're feeling well, as I come today bearing terrible news regarding your favorite chocolate chip cookie bitch")
	}
}