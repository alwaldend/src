package fp

func Compute0[T any](f func(T)) func(T) Either[T] {
	return func(t T) Either[T] {
		f(t)
		return Right(t)
	}
}

func Compute0G[T any](f func(T)) func(T) Either[T] {
	return Compute0(func(t T) {
		go f(t)
	})
}

func Compute1[T any, R any, RT Result[R]](f func(T) RT) func(T) Either[T] {
	return func(t T) Either[T] {
		if _, err := f(t).Get(); err != nil {
			return Left[T](err)
		}
		return Right(t)
	}
}
