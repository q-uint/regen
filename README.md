# regen
[![GoDoc](https://godoc.org/github.com/di-wu/regen?status.svg)](https://godoc.org/github.com/di-wu/regen)

regen is a go library that creates random strings based on regular expressions.

## usage
- `limit` is the maximum number of times `*`, `+` or `repeat`. \
  values smaller than `0` default to `10`.

```go
package main

import (
	"fmt"

	"github.com/di-wu/regen"
)

func main() {
    g, err := regen.New(`[01]{5}`)
    if err != nil {
        panic(err)
    }
    for i := 0; i < 5; i++ {
        fmt.Println(g.Generate())
    }
}
```

```
00001
01000
11101
11000
10000
```