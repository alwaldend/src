package fp

import "iter"

type Tuple[T1, T2 any] struct {
	V1 T1
	V2 T2
}

func (self Tuple[T1, T2]) Unwrap() (T1, T2) {
	return self.V1, self.V2
}

func TupleOf[T1, T2 any](v1 T1, v2 T2) Tuple[T1, T2] {
	return Tuple[T1, T2]{v1, v2}
}

func TupleIter[T1, T2 any](s iter.Seq2[T1, T2]) iter.Seq[Tuple[T1, T2]] {
	return func(yield func(Tuple[T1, T2]) bool) {
		next, stop := iter.Pull2(s)
		defer stop()
		for {
			v1, v2, ok := next()
			if !ok {
				return
			}
			if !yield(TupleOf(v1, v2)) {
				return
			}
		}
	}
}
