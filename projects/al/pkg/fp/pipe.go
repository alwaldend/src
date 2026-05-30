package fp

func Pipe(m Monad[any], ms ...func (any) Monad[any]) Monad[any] {
	for _, cur := range ms {
		m = m.FlatMap(cur)
	}
	return m
}
