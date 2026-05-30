package fp

type Box[T any] struct {
	value T
}

var _ Functor[int] = (*Box[int])(nil)
var _ Result[int] = (*Either[int])(nil)

func (self Box[T]) Map(f func(T) T) Functor[T] {
	return Box[T]{f(self.value)}
}

func (self Box[T]) Get() (T, error) {
	return self.value, nil
}
