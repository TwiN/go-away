package goaway

import (
	"testing"
)

func TestExtractProfanity(t *testing.T) {
	defaultProfanityDetector = nil
	tests := []struct {
		input             string
		expectedProfanity string
	}{
		{
			input:             "fuck this shit",
			expectedProfanity: "fuck",
		},
		{
			input:             "F   u   C  k th1$ $h!t",
			expectedProfanity: "fuck",
		},
		{
			input:             "@$$h073",
			expectedProfanity: "asshole",
		},
		{
			input:             "hello, world!",
			expectedProfanity: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			profanity := ExtractProfanity(tt.input)
			if profanity != tt.expectedProfanity {
				t.Errorf("expected '%s', got '%s'", tt.expectedProfanity, profanity)
			}
		})
	}
}

func TestProfanityDetector_Censor(t *testing.T) {
	defaultProfanityDetector = nil
	profanityDetectorWithSanitizeSpaceDisabled := NewProfanityDetector().WithSanitizeSpaces(false)
	tests := []struct {
		input                                  string
		expectedOutput                         string
		expectedOutputWithoutSpaceSanitization string
	}{
		{
			input:                                  "what the fuck",
			expectedOutput:                         "what the ****",
			expectedOutputWithoutSpaceSanitization: "what the ****",
		},
		{
			input:                                  "fuck this",
			expectedOutput:                         "**** this",
			expectedOutputWithoutSpaceSanitization: "**** this",
		},
		{
			input:                                  "one penis, two vaginas, three dicks, four sluts, five whores and a flower",
			expectedOutput:                         "one *****, two ******s, three ****s, four ****s, five *****s and a flower",
			expectedOutputWithoutSpaceSanitization: "one *****, two ******s, three ****s, four ****s, five *****s and a flower",
		},
		{
			input:                                  "Censor doesn't support sanitizing '()' into 'o', because it's two characters. Proof: c()ck. Maybe one day I'll have time to fix it.",
			expectedOutput:                         "Censor doesn't support sanitizing '()' into 'o', because it's two characters. Proof: c()ck. Maybe one day I'll have time to fix it.",
			expectedOutputWithoutSpaceSanitization: "Censor doesn't support sanitizing '()' into 'o', because it's two characters. Proof: c()ck. Maybe one day I'll have time to fix it.",
		},
		{
			input:                                  "fuck shit fuck",
			expectedOutput:                         "**** **** ****",
			expectedOutputWithoutSpaceSanitization: "**** **** ****",
		},
		{
			input:                                  "fuckfuck",
			expectedOutput:                         "********",
			expectedOutputWithoutSpaceSanitization: "********",
		},
		{
			input:                                  "fuck this shit",
			expectedOutput:                         "**** this ****",
			expectedOutputWithoutSpaceSanitization: "**** this ****",
		},
		{
			input:                                  "F   u   C  k th1$ $h!t",
			expectedOutput:                         "*   *   *  * th1$ ****",
			expectedOutputWithoutSpaceSanitization: "F   u   C  k th1$ ****",
		},
		{
			input:                                  "@$$h073",
			expectedOutput:                         "*******",
			expectedOutputWithoutSpaceSanitization: "*******",
		},
		{
			input:                                  "hello, world!",
			expectedOutput:                         "hello, world!",
			expectedOutputWithoutSpaceSanitization: "hello, world!",
		},
		{
			input:                                  "Hey asshole, are y()u an assassin? If not, fuck off.",
			expectedOutput:                         "Hey *******, are y()u an assassin? If not, **** off.",
			expectedOutputWithoutSpaceSanitization: "Hey *******, are y()u an assassin? If not, **** off.",
		},
		{
			input:                                  "I am from Scunthorpe, north Lincolnshire",
			expectedOutput:                         "I am from Scunthorpe, north Lincolnshire",
			expectedOutputWithoutSpaceSanitization: "I am from Scunthorpe, north Lincolnshire",
		},
		{
			input:                                  "He is an associate of mine",
			expectedOutput:                         "He is an associate of mine",
			expectedOutputWithoutSpaceSanitization: "He is an associate of mine",
		},
		{
			input:                                  "But the table is on fucking fire",
			expectedOutput:                         "But the table is on ****ing fire",
			expectedOutputWithoutSpaceSanitization: "But the table is on ****ing fire",
		},
		{
			input:                                  "““““““““““““But the table is on fucking fire“",
			expectedOutput:                         "““““““““““““But the table is on ****ing fire“",
			expectedOutputWithoutSpaceSanitization: "““““““““““““But the table is on ****ing fire“",
		},
		{
			input:                                  "f.u_ck this s.h-i~t",
			expectedOutput:                         "*.*_** this *.*-*~*",
			expectedOutputWithoutSpaceSanitization: "f.u_ck this s.h-i~t", // This is because special characters get replaced with a space, and because we're not sanitizing spaces...
		},
		{
			input:                                  "glass",
			expectedOutput:                         "glass",
			expectedOutputWithoutSpaceSanitization: "glass",
		},
		{
			input:                                  "ы",
			expectedOutput:                         "ы",
			expectedOutputWithoutSpaceSanitization: "ы",
		},
		{
			input:                                  "documentdocument", // false positives (https://github.com/TwiN/go-away/issues/30)
			expectedOutput:                         "documentdocument",
			expectedOutputWithoutSpaceSanitization: "documentdocument",
		},
		{
			input:                                  "dumbassdumbass", // false negatives (https://github.com/TwiN/go-away/issues/30)
			expectedOutput:                         "**************",
			expectedOutputWithoutSpaceSanitization: "**************",
		},
		{
			input:                                  "document fuck document fuck", // FIXME: This is not censored properly
			expectedOutput:                         "document **** document ****",
			expectedOutputWithoutSpaceSanitization: "document **** document ****",
		},
		{
			input:                                  "Everyone was staring, and someone muttered ‘gyat’ under their breath.",
			expectedOutput:                         "Everyone was staring, and someone muttered ‘****’ under their breath.",
			expectedOutputWithoutSpaceSanitization: "Everyone was staring, and someone muttered ‘****’ under their breath.",
		},
	}
	for _, tt := range tests {
		t.Run("default_"+tt.input, func(t *testing.T) {
			censored := Censor(tt.input)
			if censored != tt.expectedOutput {
				t.Errorf("expected '%s', got '%s'", tt.expectedOutput, censored)
			}
		})
		t.Run("no-space-sanitization_"+tt.input, func(t *testing.T) {
			censored := profanityDetectorWithSanitizeSpaceDisabled.Censor(tt.input)
			if censored != tt.expectedOutputWithoutSpaceSanitization {
				t.Errorf("expected '%s', got '%s'", tt.expectedOutputWithoutSpaceSanitization, censored)
			}
		})
	}
}

func TestNoDuplicatesBetweenProfanitiesAndFalseNegatives(t *testing.T) {
	for _, profanity := range DefaultProfanities {
		for _, falseNegative := range DefaultFalseNegatives {
			if profanity == falseNegative {
				t.Errorf("'%s' is already in 'falseNegatives', there's no need to have it in 'profanities' too", profanity)
			}
		}
	}
}

func TestBadWords(t *testing.T) {
	words := []string{"fuck", "ass", "poop", "penis", "bitch"}
	tests := []struct {
		name              string
		profanityDetector *ProfanityDetector
	}{
		{
			name:              "With Default Dictionary",
			profanityDetector: NewProfanityDetector(),
		},
		{
			name:              "With Custom Dictionary",
			profanityDetector: NewProfanityDetector().WithCustomDictionary([]string{"fuck", "ass", "poop", "penis", "bitch"}, DefaultFalsePositives, DefaultFalseNegatives),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for _, w := range words {
				if !tt.profanityDetector.IsProfane(w) {
					t.Error("Expected true, got false from word", w)
				}
				if word := tt.profanityDetector.ExtractProfanity(w); len(word) == 0 {
					t.Error("Expected true, got false from word", w)
				} else if word != w {
					t.Errorf("Expected %s, got %s", w, word)
				}
			}
		})
	}
}

func TestBadWordsWithSpaces(t *testing.T) {
	profanities := []string{"fuck", "ass", "poop", "penis", "bitch"}
	words := []string{"fu ck", "as s", "po op", "pe ni s", "bit ch"}
	tests := []struct {
		name              string
		profanityDetector *ProfanityDetector
	}{
		{
			name:              "With Default Dictionary",
			profanityDetector: NewProfanityDetector(),
		},
		{
			name:              "With Custom Dictionary",
			profanityDetector: NewProfanityDetector().WithCustomDictionary(profanities, DefaultFalsePositives, DefaultFalseNegatives),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for _, w := range words {
				if !tt.profanityDetector.WithSanitizeSpaces(true).IsProfane(w) {
					t.Error("Expected true because sanitizeSpaces is set to true, got false from word", w)
				}
				if tt.profanityDetector.WithSanitizeSpaces(false).IsProfane(w) {
					t.Error("Expected false because sanitizeSpaces is set to false, got true from word", w)
				}
			}
		})
	}
}

func TestBadWordsWithAccentedLetters(t *testing.T) {
	profanities := []string{"fuck", "ass", "poop", "penis", "bitch"}
	words := []string{"fučk", "ÄšŚ", "pÓöp", "pÉnìŚ", "bitčh"}
	tests := []struct {
		name              string
		profanityDetector *ProfanityDetector
	}{
		{
			name:              "With Default Dictionary",
			profanityDetector: NewProfanityDetector(),
		},
		{
			name:              "With Custom Dictionary",
			profanityDetector: NewProfanityDetector().WithCustomDictionary(profanities, DefaultFalsePositives, DefaultFalseNegatives),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for _, w := range words {
				if !tt.profanityDetector.WithSanitizeAccents(true).IsProfane(w) {
					t.Error("Expected true because sanitizeAccents is set to true, got false from word", w)
				}
				if tt.profanityDetector.WithSanitizeAccents(false).IsProfane(w) {
					t.Error("Expected false because sanitizeAccents is set to false, got true from word", w)
				}
			}
		})
	}
}

func TestCensorWithVerySpecialCharacters(t *testing.T) {
	profanities := []string{"крывавыa"}
	words := []string{"крывавыa"}
	expectedOutputs := []string{"********"}
	tests := []struct {
		name              string
		profanityDetector *ProfanityDetector
	}{
		{
			name:              "With Custom Dictionary",
			profanityDetector: NewProfanityDetector().WithCustomDictionary(profanities, DefaultFalsePositives, DefaultFalseNegatives),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for index, w := range words {
				if output := tt.profanityDetector.Censor(w); output != expectedOutputs[index] {
					t.Errorf("Expected %s to return %s, got %s", w, expectedOutputs[index], output)
				}
			}
		})
	}
}

func TestSentencesWithBadWords(t *testing.T) {
	profanities := []string{"fuck", "ass", "poop", "penis", "bitch"}
	sentences := []string{"What the fuck is your problem", "Go away, asshole!"}
	tests := []struct {
		name              string
		profanityDetector *ProfanityDetector
	}{
		{
			name:              "With Default Dictionary",
			profanityDetector: NewProfanityDetector(),
		},
		{
			name:              "With Custom Dictionary",
			profanityDetector: NewProfanityDetector().WithCustomDictionary(profanities, DefaultFalsePositives, DefaultFalseNegatives),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for _, s := range sentences {
				if !tt.profanityDetector.IsProfane(s) {
					t.Error("Expected true, got false from sentence", s)
				}
			}
		})
	}
}

func TestProfanityDetector_WithCustomCharacterReplacements(t *testing.T) {
	tests := []struct {
		name              string
		profanityDetector *ProfanityDetector
		sentence          string
		result            bool
	}{
		{
			name:              "With default profanity detector",
			profanityDetector: NewProfanityDetector(),
			sentence:          "5#1+",
			result:            true, // shit is a profanity
		},
		{
			name:              "With custom character replacements that has leet speak characters",
			profanityDetector: NewProfanityDetector().WithCustomCharacterReplacements(map[rune]rune{'(': 'c'}),
			sentence:          "fu(k",
			result:            true, // fuck is a profanity
		},
		{
			name:              "With custom character replacements that has leet speak characters with sanitizeLeetSpeak disabled",
			profanityDetector: NewProfanityDetector().WithCustomCharacterReplacements(map[rune]rune{'(': 'c'}).WithSanitizeLeetSpeak(false),
			sentence:          "fu(k",
			result:            false, // fuk isn't a profanity
		},
		{
			name:              "With custom character replacements that has leet speak characters with sanitizeSpecialCharacters disabled",
			profanityDetector: NewProfanityDetector().WithCustomCharacterReplacements(map[rune]rune{'(': 'c'}).WithSanitizeSpecialCharacters(false),
			sentence:          "fu(k",
			result:            false, // fu(k isn't a profanity
		},
		{
			name:              "With custom character replacements that has special characters",
			profanityDetector: NewProfanityDetector().WithCustomCharacterReplacements(map[rune]rune{'.': ' '}),
			sentence:          "f.u.c.k",
			result:            true,
		},
		{
			name:              "With custom character replacements that has special characters with sanitizeLeetSpeak disabled",
			profanityDetector: NewProfanityDetector().WithCustomCharacterReplacements(map[rune]rune{'.': ' '}).WithSanitizeLeetSpeak(false),
			sentence:          "f.u.c.k",
			result:            true, // fuck is a profanity
		},
		{
			name:              "With custom character replacements that has special characters with sanitizeSpecialCharacters disabled",
			profanityDetector: NewProfanityDetector().WithCustomCharacterReplacements(map[rune]rune{'.': ' '}).WithSanitizeSpecialCharacters(false),
			sentence:          "f.u.c.k",
			result:            false, // f.u.c.k isn't a profanity
		},
		{
			name:              "With empty character replacement mapping",
			profanityDetector: NewProfanityDetector().WithCustomCharacterReplacements(map[rune]rune{}),
			sentence:          "5#1+",
			result:            false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.profanityDetector.IsProfane(tt.sentence)
			if got != tt.result {
				t.Errorf("Expected %v, got %v from sentence %s", tt.result, got, tt.sentence)
			}
		})
	}
}

func TestSneakyBadWords(t *testing.T) {
	profanities := []string{"fuck", "ass", "poop", "penis", "bitch", "arse", "shit", "btch"}
	words := []string{"A$$", "4ss", "4s$", "a S s", "a $ s", "@$$h073", "f    u     c k", "4r5e", "5h1t", "5hit", "a55", "ar5e", "a_s_s", "b!tch", "b!+ch"}
	tests := []struct {
		name              string
		profanityDetector *ProfanityDetector
	}{
		{
			name:              "With Default Dictionary",
			profanityDetector: NewProfanityDetector(),
		},
		{
			name:              "With Custom Dictionary",
			profanityDetector: NewProfanityDetector().WithCustomDictionary(profanities, DefaultFalsePositives, DefaultFalseNegatives),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for _, w := range words {
				if !tt.profanityDetector.IsProfane(w) {
					t.Error("Expected true, got false from word", w)
				}
			}
		})
	}
}

func TestSentencesWithSneakyBadWords(t *testing.T) {
	profanities := []string{"poop", "asshole"}
	sentences := []string{
		"You smell p00p",
		"Go away, a$$h0l3!",
	}
	tests := []struct {
		name              string
		profanityDetector *ProfanityDetector
	}{
		{
			name:              "With Default Dictionary",
			profanityDetector: NewProfanityDetector(),
		},
		{
			name:              "With Custom Dictionary",
			profanityDetector: NewProfanityDetector().WithCustomDictionary(profanities, DefaultFalsePositives, DefaultFalseNegatives),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for _, s := range sentences {
				if !tt.profanityDetector.IsProfane(s) {
					t.Error("Expected true, got false from sentence", s)
				}
			}
		})
	}
}

func TestNormalWords(t *testing.T) {
	words := []string{"hello", "world", "whats", "up"}
	tests := []struct {
		name              string
		profanityDetector *ProfanityDetector
	}{
		{
			name:              "With Default Dictionary",
			profanityDetector: NewProfanityDetector(),
		},
		{
			name:              "With Custom Dictionary",
			profanityDetector: NewProfanityDetector().WithCustomDictionary(DefaultProfanities, DefaultFalsePositives, DefaultFalseNegatives),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for _, w := range words {
				if tt.profanityDetector.IsProfane(w) {
					t.Error("Expected false, got true from word", w)
				}
			}
		})
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
		"I am from Scunthorpe, north Lincolnshire",
		"He is an associate of mine",
		"Are you an assassin?",
		"But the table is on fire",
		"glass",
		"grass",
		"classic",
		"classification",
		"passion",
		"carcass",
		"cassandra",
		"just push it down the ledge", // puSH IT
		"has steph",                   // hAS Steph
		"was steph",                   // wAS Steph
		"hot water",                   // hoT WATer
		"Phoenix",                     // pHOEnix
		"systems exist",               // systemS EXist
		"saturday",                    // saTURDay
		"therapeutic",
		"press the button",
	}
	tests := []struct {
		name              string
		profanityDetector *ProfanityDetector
	}{
		{
			name:              "With Default Dictionary",
			profanityDetector: NewProfanityDetector(),
		},
		{
			name:              "With Custom Dictionary",
			profanityDetector: NewProfanityDetector().WithCustomDictionary(DefaultProfanities, DefaultFalsePositives, DefaultFalseNegatives),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for _, s := range sentences {
				if tt.profanityDetector.IsProfane(s) {
					t.Error("Expected false, got true from:", s)
				}
			}
		})
	}
}

func TestExactWord(t *testing.T) {
	sentences := []string{
		"I'm an analyst",
	}
	tests := []struct {
		name              string
		profanityDetector *ProfanityDetector
	}{
		{
			name:              "With Empty FalsePositives",
			profanityDetector: NewProfanityDetector().WithExactWord(true).WithCustomDictionary(DefaultProfanities, nil, nil),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for _, s := range sentences {
				if tt.profanityDetector.IsProfane(s) {
					t.Error("Expected false, got true from:", s)
				}
			}
		})
	}
}

func TestFalseNegatives(t *testing.T) {
	sentences := []string{
		"dumb ass", // ass -> bASS (FP) -> dumBASS (FFP)
	}
	tests := []struct {
		name              string
		profanityDetector *ProfanityDetector
	}{
		{
			name:              "With Default Dictionary",
			profanityDetector: NewProfanityDetector(),
		},
		{
			name:              "With Custom Dictionary",
			profanityDetector: NewProfanityDetector().WithCustomDictionary(DefaultProfanities, DefaultFalsePositives, DefaultFalseNegatives),
		},
		{
			name:              "With Custom Dictionary",
			profanityDetector: NewProfanityDetector().WithCustomDictionary(DefaultProfanities, DefaultFalsePositives, DefaultFalseNegatives),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for _, s := range sentences {
				if !tt.profanityDetector.IsProfane(s) {
					t.Error("Expected false, got true from:", s)
				}
			}
		})
	}
}

func TestSentencesWithFalsePositivesAndProfanities(t *testing.T) {
	sentences := []string{"You are a shitty associate", "Go away, asshole!"}
	tests := []struct {
		name              string
		profanityDetector *ProfanityDetector
	}{
		{
			name:              "With Default Dictionary",
			profanityDetector: NewProfanityDetector(),
		},
		{
			name:              "With Custom Dictionary",
			profanityDetector: NewProfanityDetector().WithCustomDictionary(DefaultProfanities, DefaultFalsePositives, DefaultFalseNegatives),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for _, s := range sentences {
				if !tt.profanityDetector.IsProfane(s) {
					t.Error("Expected true, got false from sentence", s)
				}
			}
		})
	}
}

// "The Adventures of Sherlock Holmes" by Arthur Conan Doyle is in the public domain,
// which makes it a perfect source to use as reference.
func TestSentencesFromTheAdventuresOfSherlockHolmes(t *testing.T) {
	defaultProfanityDetector = nil
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
	sanitizedString, _ := NewProfanityDetector().sanitize("What the fu_ck is y()ur pr0bl3m?", false)
	if sanitizedString != expectedString {
		t.Errorf("Expected '%s', got '%s'", expectedString, sanitizedString)
	}
}

func TestSanitizeWithoutSanitizingSpecialCharacters(t *testing.T) {
	expectedString := "whatthefu_ckisy()urproblem?"
	sanitizedString, _ := NewProfanityDetector().WithSanitizeSpecialCharacters(false).sanitize("What the fu_ck is y()ur pr0bl3m?", false)
	if sanitizedString != expectedString {
		t.Errorf("Expected '%s', got '%s'", expectedString, sanitizedString)
	}
}

func TestSanitizeWithoutSanitizingLeetSpeak(t *testing.T) {
	expectedString := "whatthefuckisyurpr0bl3m"
	sanitizedString, _ := NewProfanityDetector().WithSanitizeLeetSpeak(false).sanitize("What the fu_ck is y()ur pr0bl3m?", false)
	if sanitizedString != expectedString {
		t.Errorf("Expected '%s', got '%s'", expectedString, sanitizedString)
	}
}

func TestDefaultDriver_UTF8(t *testing.T) {
	detector := NewProfanityDetector().WithCustomDictionary(
		[]string{"anal", "あほ"}, // profanities
		[]string{"あほほ"},        // falsePositives
		[]string{"あほほし"},       // falseNegatives
	)

	unsanitizedString := "いい加減にしろ あほほし あほほ あほ anal ほ"
	expectedString := "いい加減にしろ **** あほほ ** **** ほ"

	isProfane := detector.IsProfane(unsanitizedString)
	if !isProfane {
		t.Error("Expected false, got false from sentence", unsanitizedString)
	}

	sanitizedString := detector.Censor(unsanitizedString)
	if sanitizedString != expectedString {
		t.Errorf("Expected '%s', got '%s'", expectedString, sanitizedString)
	}
}
