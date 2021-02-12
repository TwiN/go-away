package goaway

import (
	"testing"
)

func TestNoDuplicatesBetweenProfanitiesAndFalseFalsePositives(t *testing.T) {
	for _, profanity := range profanities {
		for _, falseFalsePositive := range falseNegatives {
			if profanity == falseFalsePositive {
				t.Errorf("'%s' is already in 'falseNegatives', there's no need to have it in 'profanities' too", profanity)
			}
		}
	}
}

func TestBadWords(t *testing.T) {
	words := []string{"fuck", "ass", "poop", "penis", "bitch"}
	goAway := NewProfanityDetector()
	for _, w := range words {
		if !goAway.IsProfane(w) {
			t.Error("Expected true, got false from word", w)
		}
	}
}

func TestBadWordsWithAccentedLetters(t *testing.T) {
	words := []string{"fučk", "ÄšŚ", "pÓöp", "pÉnìŚ", "bitčh"}
	for _, w := range words {
		if !NewProfanityDetector().WithSanitizeAccents(true).IsProfane(w) {
			t.Error("Expected true because sanitizeAccents is set to true, got false from word", w)
		}
		if NewProfanityDetector().WithSanitizeAccents(false).IsProfane(w) {
			t.Error("Expected false because sanitizeAccents is set to false, got true from word", w)
		}
	}
}

func TestSentencesWithBadWords(t *testing.T) {
	sentences := []string{"What the fuck is your problem", "Go away, asshole!"}
	goAway := NewProfanityDetector()
	for _, s := range sentences {
		if !goAway.IsProfane(s) {
			t.Error("Expected true, got false from sentence", s)
		}
	}
}

func TestSneakyBadWords(t *testing.T) {
	words := []string{"A$$", "4ss", "4s$", "a S s", "a $ s", "@$$h073", "f    u     c k", "4r5e", "5h1t", "5hit", "a55", "ar5e", "a_s_s", "b!tch", "b!+ch"}
	goAway := NewProfanityDetector()
	for _, w := range words {
		if !goAway.IsProfane(w) {
			t.Error("Expected true, got false from word", w)
		}
	}
}

func TestSentencesWithSneakyBadWords(t *testing.T) {
	sentences := []string{
		"You smell p00p",
		"Go away, a$$h0l3!",
	}
	goAway := NewProfanityDetector()
	for _, s := range sentences {
		if !goAway.IsProfane(s) {
			t.Error("Expected true, got false from sentence", s)
		}
	}
}

func TestNormalWords(t *testing.T) {
	words := []string{"hello", "world", "whats", "up"}
	goAway := NewProfanityDetector()
	for _, w := range words {
		if goAway.IsProfane(w) {
			t.Error("Expected false, got true from word", w)
		}
	}
}

func TestSentencesWithNoProfanities(t *testing.T) {
	sentences := []string{
		"hello, my friend",
		"what's up?",
		"do you want to play bingo?",
		"who are you?",
		"Better late than never",
		"Bite the bullet",
		"Break a leg",
		"Call it a day",
		"Be careful when you're driving",
		"How are you?",
		"Hurry up!",
		"I don't like her",
		"If you need my help, please let me know",
		"Leave a message after the beep",
		"Thank you",
		"Yes, really",
		"Call me at 9, ok?",
	}
	for _, s := range sentences {
		if IsProfane(s) {
			t.Error("Expected false, got false from sentence", s)
		}
	}
}

func TestFalsePositives(t *testing.T) {
	sentences := []string{
		"I am from Scuntthorpe, north Lincolnshire",
		"He is an associate of mine",
		"Are you an assassin?",
		"But the table is on fire",
		"glass",
		"grass",
		"classic",
		"classification",
		"passion",
		"carcass",
		"just push it down the ledge", // puSH IT
	}
	goAway := NewProfanityDetector()
	for _, s := range sentences {
		if goAway.IsProfane(s) {
			t.Error("Expected false, got true from:", s)
		}
	}
}

func TestFalseFalsePositives(t *testing.T) {
	sentences := []string{
		"dumb ass", // ass -> bASS (FP) -> dumBASS (FFP)
	}
	goAway := NewProfanityDetector()
	for _, s := range sentences {
		if !goAway.IsProfane(s) {
			t.Error("Expected false, got true from:", s)
		}
	}
}

func TestSentencesWithFalsePositivesAndProfanities(t *testing.T) {
	sentences := []string{"You are a shitty associate", "Go away, asshole!"}
	goAway := NewProfanityDetector()
	for _, s := range sentences {
		if !goAway.IsProfane(s) {
			t.Error("Expected true, got false from sentence", s)
		}
	}
}

// "The Adventures of Sherlock Holmes" by Arthur Conan Doyle is in the public domain,
// which makes it a perfect source to use as reference.
func TestSentencesFromTheAdventuresOfSherlockHolmes(t *testing.T) {
	sentences := []string{
		"I had called upon my friend, Mr. Sherlock Holmes, one day in the autumn of last year and found him in deep conversation with a very stout, florid-faced, elderly gentleman with fiery red hair.",
		"With an apology for my intrusion, I was about to withdraw when Holmes pulled me abruptly into the room and closed the door behind me.",
		"You could not possibly have come at a better time, my dear Watson, he said cordially",
		"I was afraid that you were engaged.",
		"So I am. Very much so.",
		"Then I can wait in the next room.",
		"Not at all. This gentleman, Mr. Wilson, has been my partner and helper in many of my most successful cases, and I have no doubt that he will be of the utmost use to me in yours also.",
		"The stout gentleman half rose from his chair and gave a bob of greeting, with a quick little questioning glance from his small fat-encircled eyes",
		"Try the settee, said Holmes, relapsing into his armchair and putting his fingertips together, as was his custom when in judicial moods.",
		"I know, my dear Watson, that you share my love of all that is bizarre and outside the conventions and humdrum routine of everyday life.",
		"You have shown your relish for it by the enthusiasm which has prompted you to chronicle, and, if you will excuse my saying so, somewhat to embellish so many of my own little adventures.",
		"You did, Doctor, but none the less you must come round to my view, for otherwise I shall keep on piling fact upon fact on you until your reason breaks down under them and acknowledges me to be right.",
		"Now, Mr. Jabez Wilson here has been good enough to call upon me this morning, and to begin a narrative which promises to be one of the most singular which I have listened to for some time.",
		"You have heard me remark that the strangest and most unique things are very often connected not with the larger but with the smaller crimes, and occasionally",
		"indeed, where there is room for doubt whether any positive crime has been committed.",
		"As far as I have heard it is impossible for me to say whether the present case is an instance of crime or not, but the course of events is certainly among the most singular that I have ever listened to.",
		"Perhaps, Mr. Wilson, you would have the great kindness to recommence your narrative.",
		"I ask you not merely because my friend Dr. Watson has not heard the opening part but also because the peculiar nature of the story makes me anxious to have every possible detail from your lips.",
		"As a rule, when I have heard some slight indication of the course of events, I am able to guide myself by the thousands of other similar cases which occur to my memory.",
		"In the present instance I am forced to admit that the facts are, to the best of my belief, unique.",
		"We had reached the same crowded thoroughfare in which we had found ourselves in the morning.",
		"Our cabs were dismissed, and, following the guidance of Mr. Merryweather, we passed down a narrow passage and through a side door, which he opened for us",
		"Within there was a small corridor, which ended in a very massive iron gate.",
		"We were seated at breakfast one morning, my wife and I, when the maid brought in a telegram. It was from Sherlock Holmes and ran in this way",
	}
	for _, s := range sentences {
		if IsProfane(s) {
			t.Error("Expected false, got false from sentence", s)
		}
	}
}

func TestSanitize(t *testing.T) {
	expectedString := "whatthefuckisyourproblem"
	sanitizedString := NewProfanityDetector().sanitize("What the fu_ck is y()ur pr0bl3m?")
	if sanitizedString != expectedString {
		t.Errorf("Expected '%s', got '%s'", expectedString, sanitizedString)
	}
}

func TestSanitizeWithoutSanitizingSpecialCharacters(t *testing.T) {
	expectedString := "whatthefu_ckisy()urproblem?"
	sanitizedString := NewProfanityDetector().WithSanitizeSpecialCharacters(false).sanitize("What the fu_ck is y()ur pr0bl3m?")
	if sanitizedString != expectedString {
		t.Errorf("Expected '%s', got '%s'", expectedString, sanitizedString)
	}
}

func TestSanitizeWithoutSanitizingLeetSpeak(t *testing.T) {
	expectedString := "whatthefuckisy()urpr0bl3m"
	sanitizedString := NewProfanityDetector().WithSanitizeLeetSpeak(false).sanitize("What the fu_ck is y()ur pr0bl3m?")
	if sanitizedString != expectedString {
		t.Errorf("Expected '%s', got '%s'", expectedString, sanitizedString)
	}
}
