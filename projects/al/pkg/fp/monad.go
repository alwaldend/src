package fp

type Monad[M Monad[M, T], T any] interface {
	Functor[M, T]
	FlatMap(func(T) M) M
}

type ResultMonad[M Monad[M, T], T any] interface {
	Monad[M, T]
	Result[T]
}

func FlatMap[M Monad[M, T], T any](fs ...func(v T) M) func(m M) M {
	return func(m M) M {
		for _, f := range fs {
			m = m.FlatMap(f)
		}
		return m
	}
}
