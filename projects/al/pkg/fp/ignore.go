package fp

import "iter"

func Ignore1[T1, T2, T3 any](f func(T1) (T2, T3)) func(T1) T3 {
	return func(t1 T1) T3 { _, r1 := f(t1); return r1 }
}

func IgnoreA1[T1, T2, T3 any](f func(T2) T3) func(T1, T2) T3 {
	return func(t1 T1, t2 T2) T3 { return f(t2) }
}

func IgnoreIter1[T1, T2, T3 any](f func(T1) iter.Seq2[T2, T3]) func(T1) iter.Seq[T3] {
	return func(t T1) iter.Seq[T3] {
		s := f(t)
		return func(yield func(T3) bool) {
			next, stop := iter.Pull2(s)
			defer stop()
			for {
				_, v2, ok := next()
				if !ok {
					return
				}
				if !yield(v2) {
					return
				}
			}
		}
	}
}
