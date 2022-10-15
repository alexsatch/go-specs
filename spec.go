package specs

import (
	"fmt"
)

// Spec defines an interface that enables a fluent API
// for construct specifications against T type.
//
// Each specification is immutable, and its API can be used
// to construct arbitrary trees of expressions that
// evaluate to a boolean value.
type Spec[T any] interface {
	fmt.Stringer

	// Eval evaluates a given specification against a value of type T.
	Eval(value T) bool

	// And constructs a new specification, which evaluates to true iff
	// when both this and other specification evaluate to true.
	And(other Spec[T]) Spec[T]

	// AndFunc constructs a new specification, which evaluates to true iff
	// when both this and other specification evaluate to true.
	AndFunc(other Predicate[T]) Spec[T]

	// AndNot constructs a new specification, which will be true iff
	// this specification evaluates to true and
	// other specification is false.
	AndNot(other Spec[T]) Spec[T]

	// AndNotFunc constructs a new specification, which will be true iff
	// this specification evaluates to true and
	// other specification is false.
	AndNotFunc(other Predicate[T]) Spec[T]

	// Not constructs a new specification that is a negation of this specification.
	Not() Spec[T]

	// Or constructs a new specification, which evaluates to true iff
	// when this or other specification evaluate to true.
	Or(other Spec[T]) Spec[T]

	// OrFunc constructs a new specification, which evaluates to true iff
	// when this or other specification evaluate to true.
	OrFunc(other Predicate[T]) Spec[T]

	// OrNot constructs a new specification, which evaluates to true iff
	// when this specification evaluates true or other specification evaluates to false.
	OrNot(other Spec[T]) Spec[T]

	// OrNotFunc constructs a new specification, which evaluates to true iff
	// when this specification evaluates true or other specification evaluates to false.
	OrNotFunc(other Predicate[T]) Spec[T]
}

// Predicate defines a type for a predicate on T.
type Predicate[T any] func(t T) bool

// New creates a new specification from a given predicate.
// If eval is nil, Nil specification is returned.
func New[T any](fn Predicate[T]) Spec[T] {
	if fn == nil {
		return Nil[T]()
	}

	return predicateSpec[T]{name: nameOfFunc(fn), eval: fn}
}

// NewNamed creates a new specification from a given predicate
// and gives it a name. Useful when a computed name using reflection
// doesn't fit your needs.
func NewNamed[T any](name string, fn Predicate[T]) Spec[T] {
	if fn == nil {
		return nilSpec[T]{}
	}

	if name == "" {
		panic("empty name is not permitted for specs.NewNamed()")
	}

	return predicateSpec[T]{name: name, eval: fn}
}
