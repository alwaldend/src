package fp

type IOFunc[T any] = func() (T, error)

type IO[T any] struct {
	unsafePerformIO IOFunc[T]
}

func IOOf[T any](f IOFunc[T]) IO[T] {
	return IO[T]{unsafePerformIO: f}
}

var _ Monad[IO[int], int] = (*IO[int])(nil)
var _ Functor[IO[int], int] = (*IO[int])(nil)
var _ Result[int] = (*IO[int])(nil)

func (self IO[T]) Map(f func(T) T) IO[T] {
	return IO[T]{func() (T, error) {
		res, err := self.Get()
		if err != nil {
			return res, err
		}
		return f(res), nil
	}}
}

func (self IO[T]) FlatMap(f func(T) IO[T]) IO[T] {
	return IO[T]{func() (T, error) {
		res, err := self.Get()
		if err != nil {
			return res, err
		}
		return Get(f(res))
	}}
}

func (self IO[T]) Get() (T, error) {
	return self.unsafePerformIO()
}
