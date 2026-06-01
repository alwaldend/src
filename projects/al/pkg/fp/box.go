package fp

type Box[T any] struct {
	value T
}

func BoxOf[T any](v T) Box[T] {
	return Box[T]{v}
}

var _ Functor[Box[int], int] = (*Box[int])(nil)
var _ Result[int] = (*Either[int])(nil)

func (self Box[T]) Map(f func(T) T) Box[T] {
	return Box[T]{f(self.value)}
}

func (self Box[T]) Get() (T, error) {
	return self.value, nil
}
