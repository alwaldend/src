package fp

type Maybe[T any] struct {
	Value *T
}

var _ Monad[int] = (*Maybe[int])(nil)
var _ Functor[int] = (*Maybe[int])(nil)

func Just[T any](x T) Maybe[T] {
	return Maybe[T]{&x}
}

func Nothing[T any]() Maybe[T] {
	return Maybe[T]{nil}
}

func (m Maybe[T]) Map(f func(T) T) Functor[T] {
	if m.Value == nil {
		return Nothing[T]()
	}
	return Just(f(*m.Value))
}

func (m Maybe[T]) FlatMap(f func(T) Monad[T]) Monad[T] {
	if m.Value == nil {
		return Nothing[T]()
	}
	return f(*m.Value)
}
