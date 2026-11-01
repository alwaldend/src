package fp

import "iter"

type Tuple2[T1, T2 any] struct {
	V1 T1
	V2 T2
}

func (self Tuple2[T1, T2]) Unwrap() (T1, T2) {
	return self.V1, self.V2
}

func Tuple2Of[T1, T2 any](v1 T1, v2 T2) Tuple2[T1, T2] {
	return Tuple2[T1, T2]{v1, v2}
}

func Tuple2OfE[T1, T2 any, R1 Result[T1], R2 Result[T2]](r1 R1, r2 R2) Either[Tuple2[T1, T2]] {
	v1, err := Get(r1)
	if err != nil {
		return Left[Tuple2[T1, T2]](err)
	}
	v2, err := Get(r2)
	if err != nil {
		return Left[Tuple2[T1, T2]](err)
	}
	return Right(Tuple2Of(v1, v2))
}

func Tuple2Iter[T1, T2 any](s iter.Seq2[T1, T2]) iter.Seq[Tuple2[T1, T2]] {
	return func(yield func(Tuple2[T1, T2]) bool) {
		next, stop := iter.Pull2(s)
		defer stop()
		for {
			v1, v2, ok := next()
			if !ok {
				return
			}
			if !yield(Tuple2Of(v1, v2)) {
				return
			}
		}
	}
}
