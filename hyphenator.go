/*
 * hyphenator.go
 * hyphenator
 *
 * Created by Jim Dovey on 19/07/2010.
 *
 * Copyright (c) 2010 Jim Dovey
 * All rights reserved.
 *
 * Redistribution and use in source and binary forms, with or without
 * modification, are permitted provided that the following conditions
 * are met:
 *
 * Redistributions of source code must retain the above copyright notice,
 * this list of conditions and the following disclaimer.
 *
 * Redistributions in binary form must reproduce the above copyright
 * notice, this list of conditions and the following disclaimer in the
 * documentation and/or other materials provided with the distribution.
 *
 * Neither the name of the project's author nor the names of its
 * contributors may be used to endorse or promote products derived from
 * this software without specific prior written permission.
 *
 * THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS
 * "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
 * LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS
 * FOR A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT
 * HOLDER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
 * SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED
 * TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR
 * PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF
 * LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING
 * NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS
 * SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
 *
 */

/*
	The hyphenator package implements a TeX-style hyphenation algorithm based on
	the paper by Franklin Mark Liang <http://www.tug.org/docs/liang/liang-thesis.pdf>.
*/
package hyphenator

import (
	trie "github.com/AlanQuatermain/go-trie"
	"strings"
	"io"
	"scanner"
	"os"
	"fmt"
	"utf8"
	"container/vector"
)

// The hyphenator struct itself. The nil value is a hyphenator which has not been
// initialized with any hyphenation patterns or language yet.
type Hyphenator struct {
	patterns   *trie.Trie
	exceptions map[string]string
	language   string
}

// Import hyphenation patterns from the given input stream.
func (h *Hyphenator) LoadPatterns(language string, reader io.Reader) os.Error {
	if h.language != language {
		h.patterns = nil
		h.exceptions = nil
		h.language = language
	}

	if h.patterns != nil && h.patterns.Size() != 0 {
		// looks like it's already been set up
		return nil
	}

	if reader == nil {
		return os.ErrorString("nil reader argument")
	}

	h.patterns = trie.NewTrie()
	h.exceptions = make(map[string]string, 20)

	return h.loadPatterns(reader)
}

func (h *Hyphenator) loadPatterns(reader io.Reader) os.Error {
	var s scanner.Scanner
	s.Init(reader)
	s.Mode = scanner.ScanIdents | scanner.ScanRawStrings | scanner.SkipComments

	var which string

	tok := s.Scan()
	for tok != scanner.EOF {
		switch tok {
		case scanner.Ident:
			// we handle two identifiers: 'patterns' and 'exceptions'
			switch ident := s.TokenText(); ident {
			case `patterns`, `exceptions`:
				which = ident
			default:
				return os.ErrorString(fmt.Sprintf("Unrecognized identifier '%s' at position %v",
					ident, s.Pos()))
			}
		case scanner.String, scanner.RawString:
			// trim the quotes from around the string
			tokstr := s.TokenText()
			str := tokstr[1 : len(tokstr)-1]

			switch which {
			case `patterns`:
				h.patterns.AddPatternString(str)
			case `exceptions`:
				key := strings.Replace(str, `-`, ``, -1)
				h.exceptions[key] = str
			}
		}
		tok = s.Scan()
	}
	return nil
}

func (h *Hyphenator) hyphenateWord(s, hyphen string) string {
	testStr := `.` + s + `.`
	v := make([]int, utf8.RuneCountInString(testStr))
	vIndex := 0
	for pos, _ := range testStr {
		t := testStr[pos:]
		strs, values := h.patterns.AllSubstringsAndValues(t)
		for i := 0; i < values.Len(); i++ {
			str := strs.At(i)
			val := values.At(i).(*vector.IntVector)

			diff := val.Len() - len(str)
			vs := v[vIndex-diff:]

			for i := 0; i < val.Len(); i++ {
				if val.At(i) > vs[i] {
					vs[i] = val.At(i)
				}
			}
		}
		vIndex++
	}

	var outstr string

	// trim the values for the beginning and ending dots
	markers := v[1 : len(v)-1]
	mIndex := 0
	u := make([]byte, 4)
	for _, ch := range s {
		l := utf8.EncodeRune(ch, u)
		outstr += string(u[0:l])
		// don't hyphenate between (or after) the last two characters of a string
		if mIndex < len(markers)-2 {
			// hyphens are inserted on odd values, skipped on even ones
			if markers[mIndex]%2 != 0 {
				outstr += hyphen
			}
		}
		mIndex++
	}

	return outstr
}

func (h *Hyphenator) Hyphenate(s, hyphen string) (string, bool) {
	var sc scanner.Scanner
	sc.Init(strings.NewReader(s))
	sc.Mode = scanner.ScanIdents
	sc.Whitespace = 0

	var outstr string

	tok := sc.Scan()
	for tok != scanner.EOF {
		switch tok {
		case scanner.Ident:
			// a word (or part thereof) to hyphenate
			t := sc.TokenText()

			// try the exceptions first
			exc := h.exceptions[t]
			if len(exc) != 0 {
				if hyphen != `-` {
					strings.Replace(exc, `-`, hyphen, -1)
				}
				return exc, true
			}

			// not an exception, hyphenate normally
			outstr += h.hyphenateWord(sc.TokenText(), hyphen)
		default:
			// A Unicode rune to append to the output
			p := make([]byte, utf8.UTFMax)
			l := utf8.EncodeRune(tok, p)
			outstr += string(p[0:l])
		}

		tok = sc.Scan()
	}

	return outstr, true
}
