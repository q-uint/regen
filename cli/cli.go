package main

import (
	"flag"
	"fmt"
	"time"

	"github.com/di-wu/regen"
)

func main() {
	r := flag.String("regex", "[01]{5}", "regex used to generate strings")
	s := flag.Int64("seed", time.Now().UnixNano(), "a certain seed to generate strings")
	l := flag.Int("limit", 10, "the amount of iterations of *, + and repetitions")
	n := flag.Int("n", 1, "number of regex strings to generate")
	flag.Parse()

	g, err := regen.New(*r)
	if err != nil {
		panic(err)
	}
	g.Seed(*s)
	g.Limit(*l)

	for i := 0; i < *n; i++ {
		fmt.Println(g.Generate())
	}
}
