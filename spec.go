package specs

import (
	"fmt"
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
