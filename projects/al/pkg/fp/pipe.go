package fp

type r[T any] = Result[T]

func Pipe[T1 any, R any](f1 func(T1) r[R]) func(T1) Either[R] {
	return func(t1 T1) Either[R] {
		return EitherOf(Get(f1(t1)))
	}
}

func Pipe2[T1, T2 any, R any](f1 func(T1) r[T2], f2 func(T2) r[R]) func(T1) Either[R] {
	return func(t1 T1) Either[R] {
		v, err := Get(f1(t1))
		if err != nil {
			return Left[R](err)
		}
		return Pipe(f2)(v)
	}
}

func Pipe3[T1, T2, T3 any, R any](f1 func(T1) r[T2], f2 func(T2) r[T3], f3 func(T3) r[R]) func(T1) Either[R] {
	return func(t1 T1) Either[R] {
		v, err := Get(Pipe2(f1, f2)(t1))
		if err != nil {
			return Left[R](err)
		}
		return Pipe(f3)(v)
	}
}

func Pipe4[T1, T2, T3, T4 any, R any](f1 func(T1) r[T2], f2 func(T2) r[T3], f3 func(T3) r[T4], f4 func(T4) r[R]) func(T1) Either[R] {
	return func(t1 T1) Either[R] {
		v, err := Get(Pipe3(f1, f2, f3)(t1))
		if err != nil {
			return Left[R](err)
		}
		return Pipe(f4)(v)
	}
}

func Pipe5[T1, T2, T3, T4, T5 any, R any](f1 func(T1) r[T2], f2 func(T2) r[T3], f3 func(T3) r[T4], f4 func(T4) r[T5], f5 func(T5) r[R]) func(T1) Either[R] {
	return func(t1 T1) Either[R] {
		v, err := Get(Pipe4(f1, f2, f3, f4)(t1))
		if err != nil {
			return Left[R](err)
		}
		return Pipe(f5)(v)
	}
}

func Pipe6[T1, T2, T3, T4, T5, T6 any, R any](f1 func(T1) r[T2], f2 func(T2) r[T3], f3 func(T3) r[T4], f4 func(T4) r[T5], f5 func(T5) r[T6], f6 func(T6) r[R]) func(T1) Either[R] {
	return func(t1 T1) Either[R] {
		v, err := Get(Pipe5(f1, f2, f3, f4, f5)(t1))
		if err != nil {
			return Left[R](err)
		}
		return Pipe(f6)(v)
	}
}

func Pipe7[T1, T2, T3, T4, T5, T6, T7 any, R any](f1 func(T1) r[T2], f2 func(T2) r[T3], f3 func(T3) r[T4], f4 func(T4) r[T5], f5 func(T5) r[T6], f6 func(T6) r[T7], f7 func(T7) r[R]) func(T1) Either[R] {
	return func(t1 T1) Either[R] {
		v, err := Get(Pipe6(f1, f2, f3, f4, f5, f6)(t1))
		if err != nil {
			return Left[R](err)
		}
		return Pipe(f7)(v)
	}
}

func Pipe8[T1, T2, T3, T4, T5, T6, T7, T8 any, R any](f1 func(T1) r[T2], f2 func(T2) r[T3], f3 func(T3) r[T4], f4 func(T4) r[T5], f5 func(T5) r[T6], f6 func(T6) r[T7], f7 func(T7) r[T8], f8 func(T8) r[R]) func(T1) Either[R] {
	return func(t1 T1) Either[R] {
		v, err := Get(Pipe7(f1, f2, f3, f4, f5, f6, f7)(t1))
		if err != nil {
			return Left[R](err)
		}
		return Pipe(f8)(v)
	}
}
