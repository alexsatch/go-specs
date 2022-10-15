package specs

import (
	"strings"
)

// Not creates a negation of a given spec.
func Not[T any](spec Spec[T]) Spec[T] {
	return newCompositeSpec[T]("not", combineNot[T], spec)
}

// All creates a new spec that returns true whenever all of its child specs are true.
func All[T any](specs ...Spec[T]) Spec[T] {
	return newCompositeSpec[T]("all", combineAll[T], specs...)
}

// Any creates a new spec that returns true whenever any of its child specs are true.
func Any[T any](specs ...Spec[T]) Spec[T] {
	return newCompositeSpec[T]("any", combineAny[T], specs...)
}

// newCompositeSpec returns a new compositeSpec spec.
// If list of child specs is empty, it returns a Nil spec.
func newCompositeSpec[T any](
	name string,
	combine func(childSpecs ...Spec[T]) Predicate[T],
	childSpecs ...Spec[T],
) Spec[T] {
	if len(childSpecs) == 0 {
		return Nil[T]()
	}

	return predicateSpec[T]{
		name: compositeSpecName[T](name, childSpecs...),
		eval: combine(childSpecs...),
	}
}

// compositeSpecName computes a name for a composite spec.
func compositeSpecName[T any](name string, childSpecs ...Spec[T]) string {
	sb := strings.Builder{}
	_, _ = sb.WriteString(name)

	_, _ = sb.WriteString("(")
	for i, arg := range childSpecs {
		if i != 0 {
			_, _ = sb.WriteString(", ")
		}

		_, _ = sb.WriteString(arg.String())
	}
	_, _ = sb.WriteString(")")

	return sb.String()
}

// combineNot creates a predicate that returns a negation
// of a first element of childSpecs.
func combineNot[T any](childSpecs ...Spec[T]) Predicate[T] {
	panicOnEmpty(childSpecs)

	return func(t T) bool {
		return !childSpecs[0].Eval(t)
	}
}

// combineAll creates a predicate that evaluates to true iff
// all the child specs evaluate to true.
func combineAll[T any](childSpecs ...Spec[T]) Predicate[T] {
	panicOnEmpty(childSpecs)

	return func(t T) bool {
		for _, s := range childSpecs {
			if !s.Eval(t) {
				return false
			}
		}

		return true
	}
}

// combineAll creates a predicate that evaluates to true iff
// any child spec evaluates to true.
func combineAny[T any](childSpecs ...Spec[T]) Predicate[T] {
	panicOnEmpty(childSpecs)

	return func(t T) bool {
		for _, s := range childSpecs {
			if s.Eval(t) {
				return true
			}
		}

		return false
	}
}

func panicOnEmpty[T any](childSpecs []T) {
	if len(childSpecs) == 0 {
		panic("no child specs - use nilSpec instead")
	}
}
