package fp

type Box[T any] struct {
	Value T
}

var _ Functor[int] = (*Box[int])(nil)

func (self Box[T]) Map(f func(T) T) Functor[T] {
	return Box[T]{f(self.Value)}
}
