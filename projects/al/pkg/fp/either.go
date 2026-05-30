package fp

import "fmt"

type Either[R any] struct {
	left  error
	right *R
}

var _ Monad[int] = (*Either[int])(nil)
var _ Functor[int] = (*Either[int])(nil)
var _ Result[int] = (*Either[int])(nil)

func Left[R any](err error) Either[R] {
	return Either[R]{left: err}
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
