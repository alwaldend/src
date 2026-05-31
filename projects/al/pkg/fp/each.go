package fp

import (
	"errors"
	"fmt"
	"iter"
	"runtime/debug"
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

func ForEach2[T1, T2 any](f func(v Tuple[T1, T2]) Result[any]) func(iter.Seq2[T1, T2]) EmptyEither {
	return func(s iter.Seq2[T1, T2]) EmptyEither {
		return ForEach(f)(TupleIter(s))
	}
}

func ForEach2P[T1, T2 any](f func(v Tuple[T1, T2]) Result[any]) func(iter.Seq2[T1, T2]) EmptyEither {
	return func(s iter.Seq2[T1, T2]) EmptyEither {
		return ForEachP(f)(TupleIter(s))
	}
}
