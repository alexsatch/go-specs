# go-specs
`specs` is a library that implements a [specification pattern](https://en.wikipedia.org/wiki/Specification_pattern)
against any type `T` using Go generics.

## Prerequisites
Go SDK 1.18+

## Example

```golang
package main

import (
	"github.com/alexsatch/go-specs"
)

func isPositive(x int) bool { return x > 0 }
func isNegative(x int) bool { return x < 0 }
func isZero(x int) bool     { return x == 0 }

func example1() {
	sp := specs.New(isPositive).
		OrFunc(isNegative)

	sp.Eval(1)  // true
	sp.Eval(-1) // true
	sp.Eval(0)  // false
	sp.String() // any(isPositive, isNegative)
}
```

See other examples in `examples_test.go`

