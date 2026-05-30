package fp

import "iter"

func Return[T2, T1 any](t1 T1) func(t2 T2) Monad[T1] {
	return func(t2 T2) Monad[T1] { return Right(t1) }

}

func CollectMap[T1 comparable, T2 any, T3 any](f func(k T1, v T2) Monad[T3]) func(map[T1]T2) Monad[map[T1]T3] {
	return func(m map[T1]T2) Monad[map[T1]T3] {
		res := map[T1]T3{}
		for k, v := range m {
			v2, err := Get(f(k, v))
			if err != nil {
				return Left[map[T1]T3](err)
			}
			res[k] = v2
		}
		return Right(res)
	}
}

func CollectMapV[T1 comparable, T2 any, T3 any](f func(v T2) Monad[T3]) func(map[T1]T2) Monad[map[T1]T3] {
	return CollectMap(func(k T1, v T2) Monad[T3] { return f(v) })
}

func ToAny[T1 any](v T1) Monad[any] {
	return Right[any](v)
}

func CollectSlice[T1 any, T2 any](f func(v T1) Monad[T2]) func(v []T1) Monad[[]T2] {
	return func(v []T1) Monad[[]T2] {
		res := []T2{}
		for _, v1 := range v {
			v2, err := Get(f(v1))
			if err != nil {
				return Left[[]T2](err)
			}
			res = append(res, v2)
		}
		return Right(res)
	}
}

func ForEach[T1, T2 any](seq iter.Seq[T1], t func(v T1) Monad[T2]) Monad[[]T2] {
	res := []T2{}
	for val := range seq {
		out, err := Get(t(val))
		if err != nil {
			return Left[[]T2](err)
		}
		res = append(res, out)
	}
	return Right(res)
}

func ForEach2[T1, T2, T3 any](seq iter.Seq2[T1, T2], t func(k T1, v T2) Monad[T3]) Monad[[]T3] {
	res := []T3{}
	for key, val := range seq {
		out, err := Get(t(key, val))
		if err != nil {
			return Left[[]T3](err)
		}
		res = append(res, out)
	}
	return Right(res)
}

func Compute[T1, T2 any](fs ...func(T1) Monad[T2]) func(T1) Monad[T1] {
	return func(t T1) Monad[T1] {
		for _, f := range fs {
			_, err := Get(f(t))
			if err != nil {
				return Left[T1](err)
			}
		}
		return Right(t)
	}
}

func Pipe[T any](m Monad[T], ms ...func(T) Monad[T]) Monad[T] {
	for _, cur := range ms {
		m = m.FlatMap(cur)
	}
	return m
}

func Pipe2[T1, T2 any](m Monad[T1], t1 func(T1) Monad[T2]) Monad[T2] {
	return NewIO(func() Monad[T2] {
		val, err := Get(m)
		if err != nil {
			return Left[T2](err)
		}
		return t1(val)
	})
}

func Pipe3[T1, T2, T3 any](m Monad[T1], t1 func(T1) Monad[T2], t2 func(T2) Monad[T3]) Monad[T3] {
	return NewIO(func() Monad[T3] {
		return Pipe2(Pipe2(m, t1), t2)
	})
}

func Pipe4[T1, T2, T3, T4 any](m Monad[T1], t1 func(T1) Monad[T2], t2 func(T2) Monad[T3], t3 func(T3) Monad[T4]) Monad[T4] {
	return NewIO(func() Monad[T4] {
		return Pipe2(Pipe3(m, t1, t2), t3)
	})
}

func Pipe5[T1, T2, T3, T4, T5 any](m Monad[T1], t1 func(T1) Monad[T2], t2 func(T2) Monad[T3], t3 func(T3) Monad[T4], t4 func(T4) Monad[T5]) Monad[T5] {
	return NewIO(func() Monad[T5] {
		return Pipe2(Pipe4(m, t1, t2, t3), t4)
	})
}

func Pipe6[T1, T2, T3, T4, T5, T6 any](m Monad[T1], t1 func(T1) Monad[T2], t2 func(T2) Monad[T3], t3 func(T3) Monad[T4], t4 func(T4) Monad[T5], t5 func(T5) Monad[T6]) Monad[T6] {
	return NewIO(func() Monad[T6] {
		return Pipe2(Pipe5(m, t1, t2, t3, t4), t5)
	})
}

func Pipe7[T1, T2, T3, T4, T5, T6, T7 any](m Monad[T1], t1 func(T1) Monad[T2], t2 func(T2) Monad[T3], t3 func(T3) Monad[T4], t4 func(T4) Monad[T5], t5 func(T5) Monad[T6], t6 func(T6) Monad[T7]) Monad[T7] {
	return NewIO(func() Monad[T7] {
		return Pipe2(Pipe6(m, t1, t2, t3, t4, t5), t6)
	})
}
