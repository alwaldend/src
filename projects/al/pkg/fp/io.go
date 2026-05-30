package fp

type IO[T any] struct {
	unsafePerformIO func() Monad[T]
}

func NewIO[T any](f func() Monad[T]) *IO[T] {
	return &IO[T]{f}
}

var _ Monad[int] = (*IO[int])(nil)
var _ Functor[int] = (*IO[int])(nil)
var _ Result[int] = (*IO[int])(nil)

func (self IO[T]) Map(f func(T) T) Functor[T] {
	return IO[T]{func() Monad[T] {
		res, err := self.Get()
		if err != nil {
			return Left[T](err)
		}
		return Right(f(res))
	}}
}

func (self IO[T]) Get() (T, error) {
	return Get(self.unsafePerformIO())
}

func (self IO[T]) FlatMap(f func(T) Monad[T]) Monad[T] {
	return IO[T]{func() Monad[T] {
		res, err := self.Get()
		if err != nil {
			return Left[T](err)
		}
		return f(res)
	}}
}
