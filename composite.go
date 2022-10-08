package specs

import "strings"

func Not[T any](spec Spec[T]) Spec[T]     { return newComposite[T]("not", notAggregate[T], spec) }
func All[T any](specs ...Spec[T]) Spec[T] { return newComposite[T]("all", allAggregate[T], specs...) }
func Any[T any](specs ...Spec[T]) Spec[T] { return newComposite[T]("any", anyAggregate[T], specs...) }

func newComposite[T any](
	name string,
	aggregate func(ss ...Spec[T]) func(t T) bool,
	args ...Spec[T],
) Spec[T] {
	if len(args) == 0 {
		return nilSpec[T]{}
	}

	return composite[T]{
		name:    name,
		args:    args,
		satisfy: aggregate(args...),
	}
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
	sb := strings.Builder{}
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

func notAggregate[T any](ss ...Spec[T]) func(t T) bool {
	panicIfEmpty(ss)

	return func(t T) bool {
		return !ss[0].Eval(t)
	}
}

func allAggregate[T any](ss ...Spec[T]) func(t T) bool {
	panicIfEmpty(ss)

	return func(t T) bool {
		for _, s := range ss {
			if !s.Eval(t) {
				return false
			}
		}

		return true
	}
}

func anyAggregate[T any](ss ...Spec[T]) func(t T) bool {
	panicIfEmpty(ss)

	return func(t T) bool {
		for _, s := range ss {
			if s.Eval(t) {
				return true
			}
		}

		return false
	}
}

func panicIfEmpty[T any](ss []T) {
	if len(ss) == 0 {
		panic("not supported. use nilSpec instead")
	}
}
