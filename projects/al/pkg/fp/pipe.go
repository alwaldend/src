package fp

import (
	"fmt"
	"runtime/debug"
)

type ef = func(err error) error

func Pipe[T1 any, O1 any, R1 Result[O1]](f1 func(T1) R1) func(T1) Either[O1] {
	return func(t1 T1) Either[O1] {
		return EitherOf(Get(f1(t1)))
	}
}

func PipeE[T1 any, O1 any, R1 Result[O1]](f1 func(T1) R1, fe ef) func(T1) Either[O1] {
	return func(t1 T1) Either[O1] {
		res, err := Get(f1(t1))
		if err != nil {
			return Left[O1](fe(err))
		}
		return Right(res)
	}
}

func Pipe2[T1, O1 any, R1 Result[O1], O2 any, R2 Result[O2]](f1 func(T1) R1, f2 func(O1) R2) func(T1) Either[O2] {
	return func(t1 T1) Either[O2] {
		v, err := Get(f1(t1))
		if err != nil {
			return Left[O2](err)
		}
		return Pipe(f2)(v)
	}
}

func Pipe2E[T1, O1 any, R1 Result[O1], O2 any, R2 Result[O2]](f1 func(T1) R1, f2 func(O1) R2, fe ef) func(T1) Either[O2] {
	return PipeE(Pipe2(f1, f2), fe)
}

func Pipe3[T1, O1 any, R1 Result[O1], O2 any, R2 Result[O2], O3 any, R3 Result[O3]](f1 func(T1) R1, f2 func(O1) R2, f3 func(O2) R3) func(T1) Either[O3] {
	return Pipe2(Pipe2(f1, f2), f3)
}

func Pipe3E[T1, O1 any, R1 Result[O1], O2 any, R2 Result[O2], O3 any, R3 Result[O3]](f1 func(T1) R1, f2 func(O1) R2, f3 func(O2) R3, fe ef) func(T1) Either[O3] {
	return Pipe2E(Pipe2(f1, f2), f3, fe)
}

func Pipe4[T1, O1 any, R1 Result[O1], O2 any, R2 Result[O2], O3 any, R3 Result[O3], O4 any, R4 Result[O4]](f1 func(T1) R1, f2 func(O1) R2, f3 func(O2) R3, f4 func(O3) R4) func(T1) Either[O4] {
	return Pipe3(Pipe2(f1, f2), f3, f4)
}

func Pipe4E[T1, O1 any, R1 Result[O1], O2 any, R2 Result[O2], O3 any, R3 Result[O3], O4 any, R4 Result[O4]](f1 func(T1) R1, f2 func(O1) R2, f3 func(O2) R3, f4 func(O3) R4, fe ef) func(T1) Either[O4] {
	return Pipe3E(Pipe2(f1, f2), f3, f4, fe)
}

func Pipe5[T1, O1 any, R1 Result[O1], O2 any, R2 Result[O2], O3 any, R3 Result[O3], O4 any, R4 Result[O4], O5 any, R5 Result[O5]](f1 func(T1) R1, f2 func(O1) R2, f3 func(O2) R3, f4 func(O3) R4, f5 func(O4) R5) func(T1) Either[O5] {
	return Pipe4(Pipe2(f1, f2), f3, f4, f5)
}

func Pipe5E[T1, O1 any, R1 Result[O1], O2 any, R2 Result[O2], O3 any, R3 Result[O3], O4 any, R4 Result[O4], O5 any, R5 Result[O5]](f1 func(T1) R1, f2 func(O1) R2, f3 func(O2) R3, f4 func(O3) R4, f5 func(O4) R5, fe ef) func(T1) Either[O5] {
	return Pipe4E(Pipe2(f1, f2), f3, f4, f5, fe)
}

func Pipe6[T1, O1 any, R1 Result[O1], O2 any, R2 Result[O2], O3 any, R3 Result[O3], O4 any, R4 Result[O4], O5 any, R5 Result[O5], O6 any, R6 Result[O6]](f1 func(T1) R1, f2 func(O1) R2, f3 func(O2) R3, f4 func(O3) R4, f5 func(O4) R5, f6 func(O5) R6) func(T1) Either[O6] {
	return Pipe5(Pipe2(f1, f2), f3, f4, f5, f6)
}

func Pipe6E[T1, O1 any, R1 Result[O1], O2 any, R2 Result[O2], O3 any, R3 Result[O3], O4 any, R4 Result[O4], O5 any, R5 Result[O5], O6 any, R6 Result[O6]](f1 func(T1) R1, f2 func(O1) R2, f3 func(O2) R3, f4 func(O3) R4, f5 func(O4) R5, f6 func(O5) R6, fe ef) func(T1) Either[O6] {
	return Pipe5E(Pipe2(f1, f2), f3, f4, f5, f6, fe)
}

func Pipe6EP[T1, O1 any, R1 Result[O1], O2 any, R2 Result[O2], O3 any, R3 Result[O3], O4 any, R4 Result[O4], O5 any, R5 Result[O5], O6 any, R6 Result[O6]](f1 func(T1) R1, f2 func(O1) R2, f3 func(O2) R3, f4 func(O3) R4, f5 func(O4) R5, f6 func(O5) R6, fe ef) func(T1) Either[O6] {
	return func(v T1) (res Either[O6]) {
		defer func() {
			if r := recover(); r != nil {
				res = Left[O6](fe(fmt.Errorf("panic: %s\n%s", r, string(debug.Stack()))))
			}
		}()
		return Pipe5E(Pipe2(f1, f2), f3, f4, f5, f6, fe)(v)
	}
}

func Pipe7[T1, O1 any, R1 Result[O1], O2 any, R2 Result[O2], O3 any, R3 Result[O3], O4 any, R4 Result[O4], O5 any, R5 Result[O5], O6 any, R6 Result[O6], O7 any, R7 Result[O7]](f1 func(T1) R1, f2 func(O1) R2, f3 func(O2) R3, f4 func(O3) R4, f5 func(O4) R5, f6 func(O5) R6, f7 func(O6) R7) func(T1) Either[O7] {
	return Pipe6(Pipe2(f1, f2), f3, f4, f5, f6, f7)
}

func Pipe7E[T1, O1 any, R1 Result[O1], O2 any, R2 Result[O2], O3 any, R3 Result[O3], O4 any, R4 Result[O4], O5 any, R5 Result[O5], O6 any, R6 Result[O6], O7 any, R7 Result[O7]](f1 func(T1) R1, f2 func(O1) R2, f3 func(O2) R3, f4 func(O3) R4, f5 func(O4) R5, f6 func(O5) R6, f7 func(O6) R7, fe ef) func(T1) Either[O7] {
	return Pipe6E(Pipe2(f1, f2), f3, f4, f5, f6, f7, fe)
}
