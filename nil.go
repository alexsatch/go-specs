package specs

// nilSpec implements a no-op spec that evaluates to true
// against any value of T.
//
// Its semantics are exceptional when it comes to constructing
// specification trees with it:
//  - nilSpec.Not() always returns itself;
//  - nilSpec.And*() and nilSpec.Or*() always return the (negation of) other argument;
type nilSpec[T any] struct{}

func (ns nilSpec[T]) And(other Spec[T]) Spec[T]             { return other }
func (ns nilSpec[T]) AndFunc(other Predicate[T]) Spec[T]    { return New[T](other) }
func (ns nilSpec[T]) AndNot(other Spec[T]) Spec[T]          { return Not[T](other) }
func (ns nilSpec[T]) AndNotFunc(other Predicate[T]) Spec[T] { return Not[T](New[T](other)) }
func (ns nilSpec[T]) Not() Spec[T]                          { return ns /* no conditions always eval to true */ }
func (ns nilSpec[T]) Or(other Spec[T]) Spec[T]              { return other }
func (ns nilSpec[T]) OrFunc(other Predicate[T]) Spec[T]     { return New[T](other) }
func (ns nilSpec[T]) OrNot(other Spec[T]) Spec[T]           { return Not[T](other) }
func (ns nilSpec[T]) OrNotFunc(other Predicate[T]) Spec[T]  { return Not[T](New[T](other)) }
func (ns nilSpec[T]) Eval(_ T) bool                         { return true }
func (ns nilSpec[T]) String() string                        { return "nil" }

// Nil creates a no-op specification that always evaluates to true
// on any value of T.
func Nil[T any]() Spec[T] {
	return nilSpec[T]{}
}
