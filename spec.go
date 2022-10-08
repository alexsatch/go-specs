package specs

import (
	"fmt"
	"strings"
)

type Spec[T any] interface {
	fmt.Stringer
	Eval(t T) bool

	And(other Spec[T]) Spec[T]
	AndFunc(other Func[T]) Spec[T]

	AndNot(other Spec[T]) Spec[T]
	AndNotFunc(other Func[T]) Spec[T]

	Not() Spec[T]

	Or(other Spec[T]) Spec[T]
	OrFunc(other Func[T]) Spec[T]

	OrNot(other Spec[T]) Spec[T]
	OrNotFunc(other Func[T]) Spec[T]
}

type Func[T any] func(t T) bool

type SpecFunc[T any] struct {
	name string
	fn   Func[T]
}

func (sf SpecFunc[T]) And(other Spec[T]) Spec[T]        { return All[T](sf, other) }
func (sf SpecFunc[T]) AndFunc(other Func[T]) Spec[T]    { return All[T](sf, New[T](other)) }
func (sf SpecFunc[T]) AndNot(other Spec[T]) Spec[T]     { return All[T](sf, Not[T](other)) }
func (sf SpecFunc[T]) AndNotFunc(other Func[T]) Spec[T] { return All[T](sf, Not[T](New[T](other))) }
func (sf SpecFunc[T]) Not() Spec[T]                     { return Not[T](sf) }
func (sf SpecFunc[T]) Or(other Spec[T]) Spec[T]         { return Any[T](sf, other) }
func (sf SpecFunc[T]) OrFunc(other Func[T]) Spec[T]     { return Any[T](sf, New[T](other)) }
func (sf SpecFunc[T]) OrNot(other Spec[T]) Spec[T]      { return Any[T](sf, Not[T](other)) }
func (sf SpecFunc[T]) OrNotFunc(other Func[T]) Spec[T]  { return Any[T](sf, Not[T](New[T](other))) }
func (sf SpecFunc[T]) Eval(t T) bool                    { return sf.fn(t) }
func (sf SpecFunc[T]) String() string                   { return sf.name }

type nilSpec[T any] struct{}

func (ns nilSpec[T]) And(other Spec[T]) Spec[T]        { return other }
func (ns nilSpec[T]) AndFunc(other Func[T]) Spec[T]    { return New[T](other) }
func (ns nilSpec[T]) AndNot(other Spec[T]) Spec[T]     { return Not[T](other) }
func (ns nilSpec[T]) AndNotFunc(other Func[T]) Spec[T] { return Not[T](New[T](other)) }
func (ns nilSpec[T]) Not() Spec[T]                     { return ns /* no conditions always eval to true */ }
func (ns nilSpec[T]) Or(other Spec[T]) Spec[T]         { return Any[T](ns, other) }
func (ns nilSpec[T]) OrFunc(other Func[T]) Spec[T]     { return Any[T](ns, New[T](other)) }
func (ns nilSpec[T]) OrNot(other Spec[T]) Spec[T]      { return Any[T](ns, Not[T](other)) }
func (ns nilSpec[T]) OrNotFunc(other Func[T]) Spec[T]  { return Any[T](ns, Not[T](New[T](other))) }
func (ns nilSpec[T]) Eval(_ T) bool                    { return true }
func (ns nilSpec[T]) String() string                   { return "nil" }

func New[T any](fn Func[T]) Spec[T] {
	if fn == nil {
		return nilSpec[T]{}
	}

	return SpecFunc[T]{name: nameOfFunc(fn), fn: fn}
}

type composite[T any] struct {
	name    string
	args    []Spec[T]
	satisfy func(t T) bool
}

func (c composite[T]) And(other Spec[T]) Spec[T]        { return All[T](c, other) }
func (c composite[T]) AndFunc(other Func[T]) Spec[T]    { return All[T](c, New[T](other)) }
func (c composite[T]) AndNot(other Spec[T]) Spec[T]     { return All[T](c, Not[T](other)) }
func (c composite[T]) AndNotFunc(other Func[T]) Spec[T] { return All[T](c, Not[T](New[T](other))) }
func (c composite[T]) Not() Spec[T]                     { return Not[T](c) }
func (c composite[T]) Or(other Spec[T]) Spec[T]         { return Any[T](c, other) }
func (c composite[T]) OrFunc(other Func[T]) Spec[T]     { return Any[T](c, New[T](other)) }
func (c composite[T]) OrNot(other Spec[T]) Spec[T]      { return Any[T](c, Not[T](other)) }
func (c composite[T]) OrNotFunc(other Func[T]) Spec[T]  { return Any[T](c, Not[T](New[T](other))) }

func (c composite[T]) Eval(t T) bool { return c.satisfy(t) }
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
