package fp

import "fmt"

type Either[R any] struct {
	left  error
	right *R
}

var _ Monad[Either[int], int] = (*Either[int])(nil)
var _ Functor[Either[int], int] = (*Either[int])(nil)
var _ Result[int] = (*Either[int])(nil)

type MapT[T1 comparable, T2 any] = map[T1]T2
type MapA = MapT[string, any]

func (self Either[R]) Map(f func(R) R) Either[R] {
	if self.right == nil {
		return self
	}
	return Right(f(*self.right))
}

func (self Either[R]) MapLeft(f func(err error) error) Either[R] {
	if self.left == nil {
		return self
	}
	return Left[R](f(self.left))
}

func (self Either[R]) FlatMap(f func(R) Either[R]) Either[R] {
	if self.right == nil {
		return self
	}
	return f(*self.right)
}

func (self Either[R]) Get() (R, error) {
	if self.left != nil {
		return *new(R), self.left
	}
	if self.right != nil {
		return *self.right, self.left
	}
	return *new(R), fmt.Errorf("Either is missing both left and right")
}

func ToAny(v any) Either[any] {
	return Right(v)
}

func ToAnyE(v Either[any]) Either[any] {
	return v
}

func ToAnyES(v Either[[]any]) Either[any] {
	return EitherOf[any](Get(v))
}

func EitherOf[R any](val R, err error) Either[R] {
	if err != nil {
		return Left[R](err)
	}
	return Right(val)
}

func EitherFunc2[T1, T2 any](f func(T1) (T2, error)) func(T1) Either[T2] {
	return func(t T1) Either[T2] {
		return EitherOf(f(t))
	}
}

func EitherFunc[T1 any](f func(T1) error) func(T1) EmptyEither {
	return func(t T1) EmptyEither {
		return EmptyEitherOf(f(t))
	}
}

func Left[R any](err error) Either[R] {
	return Either[R]{left: err}
}

func Right[R any](val R) Either[R] {
	return Either[R]{right: &val}
}
