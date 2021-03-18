package goaway

import (
	"testing"
)

func BenchmarkIsProfaneWhenShortStringHasNoProfanity(b *testing.B) {
	for n := 0; n < b.N; n++ {
		IsProfane("aaaaaaaaaaaaaa")
	}
	b.ReportAllocs()
}

func BenchmarkIsProfaneWhenShortStringHasProfanityAtTheStart(b *testing.B) {
	for n := 0; n < b.N; n++ {
		IsProfane("fuckaaaaaaaaaa")
	}
	b.ReportAllocs()
}

func BenchmarkIsProfaneWhenShortStringHasProfanityInTheMiddle(b *testing.B) {
	for n := 0; n < b.N; n++ {
		IsProfane("aaaaafuckaaaaa")
	}
	b.ReportAllocs()
}

func BenchmarkIsProfaneWhenShortStringHasProfanityAtTheEnd(b *testing.B) {
	for n := 0; n < b.N; n++ {
		IsProfane("aaaaaaaaaafuck")
	}
	b.ReportAllocs()
}

func BenchmarkIsProfaneWhenMediumStringHasNoProfanity(b *testing.B) {
	for n := 0; n < b.N; n++ {
		IsProfane("How are you doing today?")
	}
	b.ReportAllocs()
}

func BenchmarkIsProfaneWhenMediumStringHasProfanityAtTheStart(b *testing.B) {
	for n := 0; n < b.N; n++ {
		IsProfane("Shit, you're cute today.")
	}
	b.ReportAllocs()
}

func BenchmarkIsProfaneWhenMediumStringHasProfanityInTheMiddle(b *testing.B) {
	for n := 0; n < b.N; n++ {
		IsProfane("How are you fu ck doing?")
	}
	b.ReportAllocs()
}

func BenchmarkIsProfaneWhenMediumStringHasProfanityAtTheEnd(b *testing.B) {
	for n := 0; n < b.N; n++ {
		IsProfane("you're cute today. Fuck.")
	}
	b.ReportAllocs()
}

func BenchmarkIsProfaneWhenLongStringHasNoProfanity(b *testing.B) {
	for n := 0; n < b.N; n++ {
		IsProfane("Hello John Doe, I hope you're feeling well, as I come today bearing terrible news regarding your favorite chocolate chip cookie brand")
	}
	b.ReportAllocs()
}

func BenchmarkIsProfaneWhenLongStringHasProfanityAtTheStart(b *testing.B) {
	for n := 0; n < b.N; n++ {
		IsProfane("Fuck John Doe, I hope you're feeling well, as I come today bearing terrible news regarding your favorite chocolate chip cookie brand")
	}
	b.ReportAllocs()
}

func BenchmarkIsProfaneWhenLongStringHasProfanityInTheMiddle(b *testing.B) {
	for n := 0; n < b.N; n++ {
		IsProfane("Hello John Doe, I hope you're feeling well, as I come today bearing shitty news regarding your favorite chocolate chip cookie brand")
	}
	b.ReportAllocs()
}

func BenchmarkIsProfaneWhenLongStringHasProfanityAtTheEnd(b *testing.B) {
	for n := 0; n < b.N; n++ {
		IsProfane("Hello John Doe, I hope you're feeling well, as I come today bearing terrible news regarding your favorite chocolate chip cookie bitch")
	}
	b.ReportAllocs()
}

func BenchmarkProfanityDetector_WithSanitizeAccentsSetToFalseWhenLongStringHasProfanityAtTheStart(b *testing.B) {
	profanityDetector := NewProfanityDetector().WithSanitizeAccents(false)
	for n := 0; n < b.N; n++ {
		profanityDetector.IsProfane("Fuck John Doe, I hope you're feeling well, as I come today bearing terrible news regarding your favorite chocolate chip cookie brand")
	}
	b.ReportAllocs()
}

func BenchmarkProfanityDetector_WithSanitizeAccentsSetToFalseWhenLongStringHasProfanityInTheMiddle(b *testing.B) {
	profanityDetector := NewProfanityDetector().WithSanitizeAccents(false)
	for n := 0; n < b.N; n++ {
		profanityDetector.IsProfane("Hello John Doe, I hope you're feeling well, as I come today bearing shitty news regarding your favorite chocolate chip cookie brand")
	}
	b.ReportAllocs()
}

func BenchmarkProfanityDetector_WithSanitizeAccentsSetToFalseWhenLongStringHasProfanityAtTheEnd(b *testing.B) {
	profanityDetector := NewProfanityDetector().WithSanitizeAccents(false)
	for n := 0; n < b.N; n++ {
		profanityDetector.IsProfane("Hello John Doe, I hope you're feeling well, as I come today bearing terrible news regarding your favorite chocolate chip cookie bitch")
	}
	b.ReportAllocs()
}

func BenchmarkProfanityDetector_Sanitize(b *testing.B) {
	profanityDetector := NewProfanityDetector().WithSanitizeAccents(true).WithSanitizeSpecialCharacters(true).WithSanitizeLeetSpeak(true)
	for n := 0; n < b.N; n++ {
		profanityDetector.IsProfane("H3ll0 J0hn D0e, 1 h0p3 y0u'r3 f3eling w3ll, as 1 c0me t0d4y b34r1ng sh1tty n3w5 r3g4rd1ng y0ur fav0rite ch0c0l4t3 chip c00kie br4nd")
	}
	b.ReportAllocs()
}