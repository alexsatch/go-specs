package specs

import (
	"fmt"
	"reflect"
	"strings"
)

type Spec[T any] interface {
	fmt.Stringer
	Eval(t T) bool
}

type SpecFunc[T any] struct {
	name string
	fn   func(t T) bool
}

func (sf SpecFunc[T]) Eval(t T) bool  { return sf.fn(t) }
func (sf SpecFunc[T]) String() string { return sf.name }

func nameOfFunc[T any](fn func(t T) bool) string {
	value := reflect.ValueOf(fn)
	return value.String()
}

func New[T any](fn func(t T) bool) SpecFunc[T] {
	return SpecFunc[T]{name: nameOfFunc(fn), fn: fn}
}

type composite[T any] struct {
	name    string
	args    []Spec[T]
	satisfy func(t T) bool
}

func (c composite[T]) Eval(t T) bool {
	return c.satisfy(t)
}

func (c composite[T]) String() string {
	sb := &strings.Builder{}
	_, _ = sb.WriteString(c.name)
	sb.WriteString("(")
	for i, arg := range c.args {
		if i != 0 {
			sb.WriteString(", ")
		}

		sb.WriteString(arg.String())
	}

	sb.WriteString(")")
	return sb.String()
}

func Not[T any](spec Spec[T]) Spec[T]     { return newComposite[T]("not", notAggregate[T], spec) }
func All[T any](specs ...Spec[T]) Spec[T] { return newComposite[T]("all", allAggregate[T], specs...) }
func Any[T any](specs ...Spec[T]) Spec[T] { return newComposite[T]("any", anyAggregate[T], specs...) }

func newComposite[T any](
	name string,
	aggregate func(ss ...Spec[T]) func(t T) bool,
	args ...Spec[T],
) composite[T] {
	return composite[T]{
		name:    name,
		args:    args,
		satisfy: aggregate(args...),
	}
}

func allAggregate[T any](ss ...Spec[T]) func(t T) bool {
	return func(t T) bool {
		if len(ss) == 0 {
			return true
		}

		for _, s := range ss {
			if !s.Eval(t) {
				return false
			}
		}

		return true
	}
}

func anyAggregate[T any](ss ...Spec[T]) func(t T) bool {
	return func(t T) bool {
		if len(ss) == 0 {
			return true
		}

		for _, s := range ss {
			if s.Eval(t) {
				return true
			}
		}

		return false
	}
}

func notAggregate[T any](ss ...Spec[T]) func(t T) bool {
	return func(t T) bool {
		return !ss[0].Eval(t)
	}
}
