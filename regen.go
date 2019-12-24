package regen

import (
	"fmt"
	"math/rand"
	"regexp/syntax"
	"time"
)

const defaultLimit = 10

type r struct {
	min, max rune
}

func (r r) inRange(v rune) bool {
	return v >= r.min && v <= r.max
}

func (r r) randRune(rand *rand.Rand) rune {
	return rune(rand.Intn(int(r.max)-int(r.min)+1) + int(r.min))
}

func (r r) String() string {
	if r.min == r.max {
		return fmt.Sprintf("%s", string(r.min))
	}
	return fmt.Sprintf("%s-%s", string(r.min), string(r.max))
}

type rr []r

func (rr rr) inRange(v rune) bool {
	for _, r := range rr {
		if r.inRange(v) {
			return true
		}
	}
	return false
}

func (rr rr) randRune(rand *rand.Rand) rune {
	return rr[rand.Intn(len(rr))].randRune(rand)
}

func (rr rr) String() string {
	var str string
	for _, r := range rr {
		str += r.String()
	}
	return str
}

type Generator struct {
	re    *syntax.Regexp
	rand  *rand.Rand
	limit int
}

// Creates a new generator with given string. Returns an error if the regexp is not valid.
func New(regexp string) (*Generator, error) {
	re, err := syntax.Parse(regexp, syntax.Perl)
	if err != nil {
		return nil, err
	}

	return &Generator{
		re:    re,
		rand:  rand.New(rand.NewSource(time.Now().UnixNano())),
		limit: defaultLimit,
	}, nil
}

// Limit the amount of iterations of *, + and repetitions.
func (g *Generator) Limit(value int) {
	if value >= 0 {
		g.limit = value
	}
}

// Use a certain seed to generate strings.
func (g *Generator) Seed(seed int64) {
	g.rand = rand.New(rand.NewSource(seed))
}

// Generate a random string based on the regexp of the generator.
func (g *Generator) Generate() string {
	return g.generate(g.re)
}

func (g *Generator) generate(re *syntax.Regexp) string {
	switch re.Op {
	// matches no strings
	case syntax.OpNoMatch:
	// matches empty string
	case syntax.OpEmptyMatch,
		syntax.OpBeginLine, syntax.OpEndLine,
		syntax.OpBeginText, syntax.OpEndText:
		return ""
	// matches Runes sequence
	case syntax.OpLiteral:
		var l string
		for _, r := range re.Rune {
			l += string(r)
		}
		return l
	// matches Runes interpreted as range pair list
	case syntax.OpCharClass:
		rr := make(rr, 0)
		for i := 0; i < len(re.Rune); i += 2 {
			rr = append(rr, r{re.Rune[i], re.Rune[i+1]})
		}
		return string(rr.randRune(g.rand))
	// matches any character
	case syntax.OpAnyChar, syntax.OpAnyCharNotNL:
	// matches word (non-)boundary
	case syntax.OpWordBoundary, syntax.OpNoWordBoundary:
	// capturing subexpression with index Cap, optional name Name
	case syntax.OpCapture:
		return g.generate(re.Sub0[0])
	// matches Sub[0] zero or more times
	case syntax.OpStar:
		var l string
		for i := 0; i < g.rand.Intn(g.limit+1); i++ {
			for _, re := range re.Sub {
				l += g.generate(re)
			}
		}
		return l
	// matches Sub[0] one or more times
	case syntax.OpPlus:
		var l string
		count := g.rand.Intn(g.limit) + 1
		for i := 0; i < count; i++ {
			for _, re := range re.Sub {
				l += g.generate(re)
			}
		}
		return l
	// matches Sub[0] zero or one times
	case syntax.OpQuest:
		var l string
		for i := 0; i < g.rand.Intn(2); i++ {
			for _, re := range re.Sub {
				l += g.generate(re)
			}
		}
		return l
	// matches Sub[0] at least Min times, at most Max (Max == -1 is no limit)
	case syntax.OpRepeat:
		var l string
		if re.Max == -1 || g.limit < re.Max {
			re.Max = g.limit
		}
		var randMax int
		if re.Max > re.Min {
			randMax = g.rand.Intn(re.Max + 1)
		}
		for i := 0; i < re.Min || i < randMax; i++ {
			for _, r := range re.Sub {
				l += g.generate(r)
			}
		}
		return l
	// matches concatenation of Subs
	case syntax.OpConcat:
		var l string
		for _, re := range re.Sub {
			l += g.generate(re)
		}
		return l
	// matches alternation of Subs
	case syntax.OpAlternate:
		return g.generate(re.Sub[g.rand.Intn(len(re.Sub))])
	}
	return ""
}
