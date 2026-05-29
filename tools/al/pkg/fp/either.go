package fp

type Either[L, R any] struct {
	Left  *L
	Right *R
}

var _ Monad[int] = (*Either[error, int])(nil)
var _ Functor[int] = (*Either[error, int])(nil)

func Left[L, R any](x L) Either[L, R] {
	return Either[L, R]{Left: &x}
}

func Right[L, R any](x R) Either[L, R] {
	return Either[L, R]{Right: &x}
}

func (self Either[L, R]) Map(f func(R) R) Functor[R] {
	if self.Right == nil {
		return self
	}
	return Right[L](f(*self.Right))
}

func (self Either[L, R]) FlatMap(f func(R) Monad[R]) Monad[R] {
	if self.Right == nil {
		return self
	}
	return f(*self.Right)
}
