package hyphenator

import (
	"testing"
	"fmt"
	"os"
)
/*
const testStr = `Out of a list with hyphenated words he computed patterns with integer values.
Odd numbers are marking hyphenation points. With this patterns a
program can compute possible hyphen points of any given word.`

const hyphStr = `Out of a list with hy-phen-at-ed words he com-put-ed pat-terns with in-te-ger val-ues.
Odd num-bers are mark-ing hy-phen-ation points. With this pat-terns a
pro-gram can com-pute pos-si-ble hy-phen points of any giv-en word.`
*/

// technically this string contains an em-dash character, but the scanner.Scanner barfs on that for some
// reason, producing error glyphs in the output.  It also parses it twice, which is super-annoying.  For
// this reason I've replaced it with a double-hyphen sequence, like many ASCII-limited people before me.
const testStr = `Go is a new language. Although it borrows ideas from existing languages, it has unusual properties that make effective Go programs different in character from programs written in its relatives. A straightforward translation of a C++ or Java program into Go is unlikely to produce a satisfactory result--Java programs are written in Java, not Go. On the other hand, thinking about the problem from a Go perspective could produce a successful but quite different program. In other words, to write Go well, it's important to understand its properties and idioms. It's also important to know the established conventions for programming in Go, such as naming, formatting, program construction, and so on, so that programs you write will be easy for other Go programmers to understand.

This document gives tips for writing clear, idiomatic Go code. It augments the language specification and the tutorial, both of which you should read first.

Examples
The Go package sources are intended to serve not only as the core library but also as examples of how to use the language. If you have a question about how to approach a problem or how something might be implemented, they can provide answers, ideas and background.

Formatting
Formatting issues are the most contentious but the least consequential. People can adapt to different formatting styles but it's better if they don't have to, and less time is devoted to the topic if everyone adheres to the same style. The problem is how to approach this Utopia without a long prescriptive style guide.

With Go we take an unusual approach and let the machine take care of most formatting issues. A program, gofmt, reads a Go program and emits the source in a standard style of indentation and vertical alignment, retaining and if necessary reformatting comments. If you want to know how to handle some new layout situation, run gofmt; if the answer doesn't seem right, fix the program (or file a bug), don't work around it.

As an example, there's no need to spend time lining up the comments on the fields of a structure. Gofmt will do that for you.`

const hyphStr = `Go is a new lan-guage. Although it bor-rows ideas from ex-ist-ing lan-guages, it has un-usu-al prop-er-ties that make ef-fec-tive Go pro-grams d-if-fer-ent in char-ac-ter from pro-grams writ-ten in its rel-a-tives. A s-traight-for-ward trans-la-tion of a C++ or Ja-va pro-gram in-to Go is un-like-ly to pro-duce a sat-is-fac-to-ry re-sult--Ja-va pro-grams are writ-ten in Ja-va, not Go. On the oth-er hand, think-ing about the prob-lem from a Go per-spec-tive could pro-duce a suc-cess-ful but quite d-if-fer-ent pro-gram. In oth-er words, to write Go well, it's im-por-tant to un-der-stand its prop-er-ties and id-ioms. It's al-so im-por-tant to know the es-tab-lished con-ven-tions for pro-gram-ming in Go, such as nam-ing, for-mat-ting, pro-gram con-struc-tion, and so on, so that pro-grams y-ou write will be easy for oth-er Go pro-gram-mers to un-der-stand.

This doc-u-ment gives tips for writ-ing clear, id-iomat-ic Go code. It aug-ments the lan-guage spec-i-fi-ca-tion and the tu-to-r-i-al, both of which y-ou should read first.

Ex-am-ples
The Go pack-age sources are in-tend-ed to serve not on-ly as the core li-brary but al-so as ex-am-ples of how to use the lan-guage. If y-ou have a ques-tion about how to ap-proach a prob-lem or how some-thing might be im-ple-ment-ed, they can pro-vide an-swers, ideas and back-ground.

For-mat-ting
For-mat-ting is-sues are the most con-tentious but the least con-se-quen-tial. Peo-ple can adapt to d-if-fer-ent for-mat-ting styles but it's bet-ter if they don't have to, and less time is de-vot-ed to the top-ic if every-one ad-heres to the same style. The prob-lem is how to ap-proach this U-topia with-out a long pre-scrip-tive style guide.

With Go we take an un-usu-al ap-proach and let the ma-chine take care of most for-mat-ting is-sues. A pro-gram, gofmt, reads a Go pro-gram and emits the source in a s-tan-dard style of in-den-ta-tion and ver-ti-cal align-ment, re-tain-ing and if nec-es-sary re-for-mat-ting com-ments. If y-ou want to know how to han-dle some new lay-out sit-u-a-tion, run gofmt; if the an-swer doesn't seem right, fix the pro-gram (or file a bug), don't work around it.

As an ex-am-ple, there's no need to spend time lin-ing up the com-ments on the fields of a struc-ture. Gofmt will do that for y-ou.`

func buildHyphenator() *Hyphenator {
	h := new(Hyphenator)
	f, err := os.Open("patterns-en", 0666, os.O_RDONLY)
	if err != nil {
		fmt.Println("osOpen():", err)
		os.Exit(1)
	}

	//fmt.Println("Loading patterns...")
	h.LoadPatterns("en", f)
	f.Close()
	//fmt.Println("...done.\n\n")

	//fmt.Println("Size of compiled hyphenation trie:", h.patterns.Size())
	//fmt.Println("Size of compiled exception list:", len(h.exceptions))
	return h
}

func TestHyphenator(t *testing.T) {
	h := buildHyphenator()
	hyphenated, ok := h.Hyphenate(testStr, `-`)
	if !ok {
		t.Fail()
	}
	//fmt.Println(hyphenated)
	if hyphenated != hyphStr {
		t.Fail()
	}
}

func BenchmarkHyphenation(b *testing.B) {
	b.StopTimer()
	h := buildHyphenator()

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		h.Hyphenate(testStr, `-`)
	}
}

func BenchmarkHTMLHyphenation(b *testing.B) {
	b.StopTimer()
	h := buildHyphenator()

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		h.Hyphenate(testStr, `&shy;`)
	}
}
