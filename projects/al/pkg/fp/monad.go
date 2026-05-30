package fp

type Monad[T any] interface {
	Functor[T]
	FlatMap(func(T) Monad[T]) Monad[T]
}

type EmptyMonad = Monad[struct{}]
