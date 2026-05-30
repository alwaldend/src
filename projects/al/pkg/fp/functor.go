package fp

type Functor[T any] interface {
	Map(func(T) T) Functor[T]
}

func Map[T any](v Functor[T], f func(T) T) Functor[T] {
	return v.Map(f)
}
