package fp

import "fmt"

func Switch[T comparable, R any](v1 T, r1 R) func(v T) Either[R] {
	return func(v T) Either[R] {
		switch v {
		case v1:
			return Right(r1)
		default:
			return Left[R](fmt.Errorf("missing case for %v: ", v))
		}
	}
}

func Switch2[T comparable, R any](v1 T, r1 R, v2 T, r2 R) func(v T) Either[R] {
	return func(v T) Either[R] {
		switch v {
		case v1:
			return Right(r1)
		case v2:
			return Right(r2)
		default:
			return Left[R](fmt.Errorf("missing case for %v: ", v))
		}
	}
}

func Switch3[T comparable, R any](v1 T, r1 R, v2 T, r2 R, v3 T, r3 R) func(v T) Either[R] {
	return func(v T) Either[R] {
		switch v {
		case v1:
			return Right(r1)
		case v2:
			return Right(r2)
		case v3:
			return Right(r3)
		default:
			return Left[R](fmt.Errorf("missing case for %v: ", v))
		}
	}
}

func Switch4[T comparable, R any](v1 T, r1 R, v2 T, r2 R, v3 T, r3 R, v4 T, r4 R) func(v T) Either[R] {
	return func(v T) Either[R] {
		switch v {
		case v1:
			return Right(r1)
		case v2:
			return Right(r2)
		case v3:
			return Right(r3)
		case v4:
			return Right(r4)
		default:
			return Left[R](fmt.Errorf("missing case for %v: ", v))
		}
	}
}
