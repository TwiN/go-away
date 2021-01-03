![go-away](/.github/assets/go-away.png)

# go-away

![build](https://github.com/TwinProduction/go-away/workflows/build/badge.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/TwinProduction/go-away)](https://goreportcard.com/report/github.com/TwinProduction/go-away)
[![codecov](https://codecov.io/gh/TwinProduction/go-away/branch/master/graph/badge.svg)](https://codecov.io/gh/TwinProduction/go-away)
[![Go Reference](https://pkg.go.dev/badge/github.com/TwinProduction/go-away.svg)](https://pkg.go.dev/github.com/TwinProduction/go-away)

go-away is a stand-alone, lightweight library for detecting profanities in Go.

This library must remain **extremely** easy to use. Its original intent of not adding overhead will always remain.


## Installation

```
go get -u github.com/TwinProduction/go-away
```


## Usage

```go
import (
	"github.com/TwinProduction/go-away"
)

goaway.IsProfane("fuck this shit")         // returns true
goaway.IsProfane("F   u   C  k th1$ $h!t") // returns true
goaway.IsProfane("@$$h073")                // returns true
goaway.IsProfane("hello, world!")          // returns false
```

By default, `IsProfane` uses a default profanity detector, but if you'd like to disable leet speak,
numerical character or special character sanitization, you have to create a ProfanityDetector instead:

```go
profanityDetector := goaway.NewProfanityDetector().WithSanitizeLeetSpeak(false).WithSanitizeSpecialCharacters(false).WithSanitizeAccents(false)
profanityDetector.IsProfane("b!tch") // returns false because we're not sanitizing special characters
```


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

