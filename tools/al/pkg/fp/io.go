package fp

type IO[T any] struct {
	UnsafePerformIO func() T
}

func NewIO[T any](f func() T) *IO[T] {
	return &IO[T]{UnsafePerformIO: f}
}

func (self IO[T]) Map(f func(T) T) Functor[T] {
	return IO[T]{func() T {
		return f(self.UnsafePerformIO())
	}}
}

func (self IO[T]) FlatMap(f func(T) Monad[T]) Monad[T] {
	return IO[T]{func() T {
		return f(self.UnsafePerformIO()).(IO[T]).UnsafePerformIO()
	}}
}
