# go-away

go-away is a stand-alone, lightweight library for detecting profanities in Go.

This library must remain **extremely** easy to use. Its original intent of not adding overhead will always remain.

## Installation

```
go get github.com/TwinProduction/go-away
```

## Usage

```golang
import (
	"github.com/TwinProduction/go-away"
)

goaway.IsProfane("fuck this shit") // returns true
goaway.IsProfane("F   u   C  k th1$ $h!t") // returns true
goaway.IsProfane("@$$h073") // returns true
goaway.IsProfane("hello, world!") // returns false
```

## In the background

While using a giant regex query to handle everything would be a way of doing it, as more words 
are added to the list of profanities, that would slow down the filtering considerably.

Instead, the following steps are taken before checking for profanities in a string:

- Numbers are replaced to their letter counterparts (e.g. 1 -> L, 4 -> A, etc)
- Special characters are replaced to their letter equivalent (e.g. @ -> A, ! -> i)
- The resulting string has all of its spaces removed, to prevent `w  ords  lik e   tha   t`
- The string has all of its characters converted to lowercase
- [**TODO**] All non-transformed special characters are removed, to prevent `w__ords li.ke tha--t`
- [**TODO**] All words that have the same character repeated more than twice in a row are removed (e.g. `poooop -> poop`)
    - NOTE: This is obviously not a perfect approach, as words like `fuuck` wouldn't be detected, but it's better than nothing.
    

The upside of this method is that we only need to add base bad words, and not all tenses of said bad word.

e.g. the `fuck` entry would support `fucker`, `fucking`, etc)

The downside is that words like `assassin`, which contain `ass`, would also be filtered as profane.
So in the future, a list of false positives would have to be added.
