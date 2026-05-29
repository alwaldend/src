package fp

type Functor[T any] interface {
	Map(func(T) T) Functor[T]
}
