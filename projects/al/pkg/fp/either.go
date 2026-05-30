package fp

import "fmt"

type Either[R any] struct {
	left  error
	right *R
}

var _ Monad[int] = (*Either[int])(nil)
var _ Functor[int] = (*Either[int])(nil)
var _ Result[int] = (*Either[int])(nil)

type EmptyEither = Either[struct{}]

func NewEmptyEither(err error) EmptyEither {
	if err == nil {
		return Right(struct{}{})
	}
	return Left[struct{}](err)
}

func NewEither[R any](val R, err error) Either[R] {
	if err != nil {
		return Left[R](err)
	}
	return Right(val)
}

func EitherFunc2[T1, T2 any](f func(T1) (T2, error)) func(T1) Monad[T2] {
	return func(t T1) Monad[T2] {
		res, err := f(t)
		if err != nil {
			return Left[T2](err)
		}
		return Right(res)
	}
}

func EitherFunc[T1 any](f func(T1) error) func(T1) EmptyMonad {
	return func(t T1) EmptyMonad {
		err := f(t)
		if err != nil {
			return EmptyLeft(err)
		}
		return EmptyRight()
	}
}

func EmptyLeft(err error) Either[struct{}] {
	return Left[struct{}](err)
}

func Left[R any](err error) Either[R] {
	return Either[R]{left: err}
}

func EmptyRight() Either[struct{}] {
	return Right(struct{}{})
}

func Right[R any](val R) Either[R] {
	return Either[R]{right: &val}
}

func (self Either[R]) Map(f func(R) R) Functor[R] {
	if self.right == nil {
		return self
	}
	return Right(f(*self.right))
}

func (self Either[R]) FlatMap(f func(R) Monad[R]) Monad[R] {
	if self.right == nil {
		return self
	}
	return f(*self.right)
}

func (self Either[R]) Get() (R, error) {
	if self.left != nil {
		var r R
		return r, self.left
	}
	if self.right != nil {
		return *self.right, self.left
	}
	var r R
	return r, fmt.Errorf("Either is missing both left and right")
}
