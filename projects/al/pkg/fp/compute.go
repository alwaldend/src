package fp

func Compute0[T any](fs ...func(T)) func(T) Either[T] {
	return func(t T) Either[T] {
		for _, f := range fs {
			f(t)
		}
		return Right(t)
	}
}

func Compute1[T any, R any, RT Result[R]](fs ...func(T) RT) func(T) Either[T] {
	return func(t T) Either[T] {
		for _, f := range fs {
			_, err := f(t).Get()
			if err != nil {
				return Left[T](err)
			}
		}
		return Right(t)
	}
}
