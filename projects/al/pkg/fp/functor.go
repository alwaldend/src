package fp

type Functor[F Functor[F, T], T any] interface {
	Map(func(T) T) F
}

type ResultFunctor[F Functor[F, T], T any] interface {
	Functor[F, T]
	Result[T]
}

func Map[F Functor[F, T], T any](fs ...func(T) T) func(f F) F {
	return func(f F) F {
		for _, fn := range fs {
			f = f.Map(fn)
		}
		return f
	}
}
