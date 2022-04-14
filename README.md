![go-away](/.github/assets/go-away.png)

# go-away
![test](https://github.com/TwiN/go-away/workflows/test/badge.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/TwiN/go-away)](https://goreportcard.com/report/github.com/TwiN/go-away)
[![codecov](https://codecov.io/gh/TwiN/go-away/branch/master/graph/badge.svg)](https://codecov.io/gh/TwiN/go-away)
[![Go Reference](https://pkg.go.dev/badge/github.com/TwiN/go-away.svg)](https://pkg.go.dev/github.com/TwiN/go-away)
[![Follow TwiN](https://img.shields.io/github/followers/TwiN?label=Follow&style=social)](https://github.com/TwiN)

go-away is a stand-alone, lightweight library for detecting and censoring profanities in Go.

This library must remain **extremely** easy to use. Its original intent of not adding overhead will always remain.


## Installation
```console
go get -u github.com/TwiN/go-away
```


## Usage
```go
package main

import (
    "github.com/TwiN/go-away"
)

func main() {
    goaway.IsProfane("fuck this shit")                // returns true
    goaway.ExtractProfanity("fuck this shit")         // returns "fuck"
    goaway.Censor("fuck this shit")                   // returns "**** this ****"
    
    goaway.IsProfane("F   u   C  k th1$ $h!t")        // returns true
    goaway.ExtractProfanity("F   u   C  k th1$ $h!t") // returns "fuck"
    goaway.Censor("F   u   C  k th1$ $h!t")           // returns "*   *   *  * th1$ ****"
    
    goaway.IsProfane("@$$h073")                       // returns true
    goaway.ExtractProfanity("@$$h073")                // returns "asshole"
    goaway.Censor("@$$h073")                          // returns "*******"
    
    goaway.IsProfane("hello, world!")                 // returns false
    goaway.ExtractProfanity("hello, world!")          // returns ""
    goaway.Censor("hello, world!")                    // returns "hello, world!"
}
```

Calling `goaway.IsProfane(s)`, `goaway.ExtractProfanity(s)` or `goaway.Censor(s)` will use the default profanity detector,
but if you'd like to disable leet speak, numerical character or special character sanitization, you have to create a
ProfanityDetector instead:
```go
profanityDetector := goaway.NewProfanityDetector().WithSanitizeLeetSpeak(false).WithSanitizeSpecialCharacters(false).WithSanitizeAccents(false)
profanityDetector.IsProfane("b!tch") // returns false because we're not sanitizing special characters
```

By default, the `NewProfanityDetector` constructor uses the default dictionaries for profanities, false positives and false negatives.
These dictionaries are exposed as `goaway.DefaultProfanities`, `goaway.DefaultFalsePositives` and `goaway.DefaultFalseNegatives` respectively.

If you need to load a different dictionary, you could create a new instance of `ProfanityDetector` on this way:
```go
profanities    := []string{"ass"}
falsePositives := []string{"bass"}
falseNegatives := []string{"dumbass"}

profanityDetector := goaway.NewProfanityDetector().WithCustomDictionary(profanities, falsePositives, falseNegatives)
```

You may also specify custom character replacements using `WithCustomCharacterReplacements` on a `ProfanityDetector`.
By default, this is set to `goaway.DefaultCharacterReplacements`.

Note that all character replacements with a value of `' '` are considered as special characters while all characters
with a value that is not `' '` are considered to be leetspeak characters. This means that using 
`profanityDetector.WithSanitizeSpecialCharacters(bool)` and `profanityDetector.WithSanitizeLeetSpeak(bool)` will let you
toggle which character replacements are executed during the sanitization process.


## In the background
While using a giant regex query to handle everything would be a way of doing it, as more words 
are added to the list of profanities, that would slow down the filtering considerably.

Instead, the following steps are taken before checking for profanities in a string:

- Numbers are replaced to their letter counterparts (e.g. 1 -> L, 4 -> A, etc)
- Special characters are replaced to their letter equivalent (e.g. @ -> A, ! -> i)
- The resulting string has all of its spaces removed to prevent `w  ords  lik e   tha   t`
- The resulting string has all of its characters converted to lowercase
- The resulting string has all words deemed as false positives (e.g. `assassin`) removed

In the future, the following additional steps could also be considered:
- All non-transformed special characters are removed to prevent `s~tring li~ke tha~~t`
- All words that have the same character repeated more than twice in a row are removed (e.g. `poooop -> poop`)
  - NOTE: This is obviously not a perfect approach, as words like `fuuck` wouldn't be detected, but it's better than nothing.
  - The upside of this method is that we only need to add base bad words, and not all tenses of said bad word. (e.g. the `fuck` entry would support `fucker`, `fucking`, etc.)
