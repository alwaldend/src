package fp

import (
	"errors"
	"fmt"
	"iter"
	"runtime/debug"
	"slices"
	"sync"
)

func ForEach[T any](f func(T) Result[any]) func(iter.Seq[T]) EmptyEither {
	return func(s iter.Seq[T]) EmptyEither {
		for v := range s {
			if _, err := f(v).Get(); err != nil {
				return EmptyLeft(err)
			}
		}
		return EmptyRight()
	}
}

func CollectEach[T1, T2 any](f func(T1) Result[T2]) func(iter.Seq[T1]) Either[[]T2] {
	return func(s iter.Seq[T1]) Either[[]T2] {
		res := []T2{}
		for v := range s {
			v2, err := Get(f(v))
			if err != nil {
				return Left[[]T2](err)
			}
			res = append(res, v2)
		}
		return Right(res)
	}
}

func CollectEachS[T1, T2 any](f func(T1) Result[T2]) func([]T1) Either[[]T2] {
	return func(s []T1) Either[[]T2] {
		it := IgnoreIter1[[]T1](slices.All)(s)
		return CollectEach(f)(it)
	}
}

func CollectEachM[M ~map[T1]T2, T1 comparable, T2, T3 any](f func(T1, T2) Result[T3]) func(M) Either[map[T1]T3] {
	return func(m M) Either[map[T1]T3] {
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

func CollectEachMV[M ~map[T1]T2, T1 comparable, T2, T3 any](f func(T2) Result[T3]) func(M) Either[map[T1]T3] {
	return CollectEachM[M](IgnoreA1[T1](f))
}

func ForEachP[T any](f func(v T) Result[any]) func(iter.Seq[T]) EmptyEither {
	return func(s iter.Seq[T]) EmptyEither {
		var wg sync.WaitGroup
		errs := make(chan error)
		errsConsumer := make(chan error)
		f2 := func(v T) Result[any] {
			wg.Go(func() {
				defer func() {
					if r := recover(); r != nil {
						errs <- fmt.Errorf("panic: %s\n%s", r, debug.Stack())
					}
				}()
				if _, err := f(v).Get(); err != nil {
					errs <- err
				}
			})
			return Right[any](struct{}{})
		}
		_, err := ForEach(f2)(s).Get()
		go func() {
			for err2 := range errs {
				errsConsumer <- err2
			}
			close(errsConsumer)
		}()
		wg.Wait()
		close(errs)
		for e := range errs {
			err = errors.Join(err, e)
		}
		if err != nil {
			return EmptyLeft(err)
		}
		return EmptyRight()
	}
}

func ForEach2[T1, T2 any](f func(v Tuple2[T1, T2]) Result[any]) func(iter.Seq2[T1, T2]) EmptyEither {
	return func(s iter.Seq2[T1, T2]) EmptyEither {
		return ForEach(f)(Tuple2Iter(s))
	}
}

func ForEach2P[T1, T2 any](f func(v Tuple2[T1, T2]) Result[any]) func(iter.Seq2[T1, T2]) EmptyEither {
	return func(s iter.Seq2[T1, T2]) EmptyEither {
		return ForEachP(f)(Tuple2Iter(s))
	}
}
