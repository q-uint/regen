package regen

import (
	"math/rand"
	"regexp"
	"testing"
	"time"
)

var (
	// [0-9]
	digits = r{48, 57}
	// [a-z]
	lowercase = r{97, 122}
	// [A-Z]
	uppercase = r{65, 90}
	// [A-Za-z]
	alphabetic = rr{lowercase, uppercase}
	// [0-9A-Za-z]
	alphanumeric = append(alphabetic, digits)
	// [\t ]
	blank = rr{{9, 9}, {32, 32}}
	// [\t\n\v\f\r ]
	whitespace = append(blank, r{10, 13})
	// [!-/:-@[-`{-~]
	punctuation = rr{{33, 47}, {58, 64}, {91, 96}, {123, 126}}
	// [!-~]
	graphical = r{33, 126}
)

const iterations = 100

func TestRandRune(t *testing.T) {
	seed := rand.New(rand.NewSource(time.Now().UnixNano()))

	for _, test := range []r{
		digits, lowercase, uppercase, graphical,
	} {
		for i := 0; i < iterations; i++ {
			r := test.randRune(seed)
			if !test.inRange(r) {
				t.Error(test, r)
			}
		}

		for r := test.min; r <= test.max; r++ {
			if !test.inRange(r) {
				t.Error(test, r)
			}
		}
	}

	for _, test := range []rr{
		alphabetic, alphanumeric, blank, whitespace, punctuation,
	} {
		for i := 0; i < iterations; i++ {
			r := test.randRune(seed)
			if !test.inRange(r) {
				t.Error(test, r)
			}
		}
	}
}

func TestGenerate(t *testing.T) {
	for _, test := range []string{
		`[abc]`,    // a single character of: a, b or c
		`[^abc]`,   // any single character except: a, b, or c
		`[a-z]`,    // any single character in the range a-z
		`[a-zA-Z]`, // any single character in the range a-z or A-Z
		`(a|b)`,    // a or b
		`a?`,       // zero or one of a
		`a*`,       // zero or more of a
		`a+`,       // one or more of a
		`a{3}`,     // exactly 3 of a
		`a{3,}`,    // 3 or more of a
		`a{3,6}`,   // between 3 and 6 of a
	} {
		g, _ := New(test)
		for i := 0; i < iterations; i++ {
			re := g.Generate()

			if !regexp.MustCompile(test).MatchString(re) {
				t.Error(test, re)
			}
		}
	}
}

func TestSeed(t *testing.T) {
	test := `[a-zA-Z]{10}`
	g1, _ := New(test)
	g2, _ := New(test)
	for i := 0; i < iterations; i++ {
		seed := rand.Int63()
		g1.Seed(seed)
		g2.Seed(seed)
		re1 := g1.Generate()
		re2 := g2.Generate()
		if re1 != re2 {
			t.Error(seed, re1, re2)
		}
	}
}

func TestLimit(t *testing.T) {
	for _, limit := range []int{
		0, 1, 10, -1,
	} {
		g, _ := New(`[a-zA-Z]*`)
		g.Limit(limit)
		for i := 0; i < iterations; i++ {
			re := g.Generate()
			if len(re) > 10 {
				t.Error(limit, len(re), re)
			}
		}
	}
}
