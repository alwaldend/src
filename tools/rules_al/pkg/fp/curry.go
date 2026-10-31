package fp

func Curry1[T1, T2 any](v1 T1) func(f1 func(T1) T2) func() T2 {
	return func(f1 func(T1) T2) func() T2 {
		return func() T2 {
			return f1(v1)
		}
	}
}
