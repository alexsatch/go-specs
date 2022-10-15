package specs

// predicateSpec enables a fluent API over specification predicates.
type predicateSpec[T any] struct {
	name string
	eval Predicate[T]
}

func (ps predicateSpec[T]) And(other Spec[T]) Spec[T]          { return All[T](ps, other) }
func (ps predicateSpec[T]) AndFunc(other Predicate[T]) Spec[T] { return All[T](ps, New[T](other)) }
func (ps predicateSpec[T]) AndNot(other Spec[T]) Spec[T]       { return All[T](ps, Not[T](other)) }
func (ps predicateSpec[T]) AndNotFunc(other Predicate[T]) Spec[T] {
	return All[T](ps, Not[T](New[T](other)))
}
func (ps predicateSpec[T]) Not() Spec[T]                      { return Not[T](ps) }
func (ps predicateSpec[T]) Or(other Spec[T]) Spec[T]          { return Any[T](ps, other) }
func (ps predicateSpec[T]) OrFunc(other Predicate[T]) Spec[T] { return Any[T](ps, New[T](other)) }
func (ps predicateSpec[T]) OrNot(other Spec[T]) Spec[T]       { return Any[T](ps, Not[T](other)) }
func (ps predicateSpec[T]) OrNotFunc(other Predicate[T]) Spec[T] {
	return Any[T](ps, Not[T](New[T](other)))
}
func (ps predicateSpec[T]) Eval(t T) bool  { return ps.eval(t) }
func (ps predicateSpec[T]) String() string { return ps.name }
