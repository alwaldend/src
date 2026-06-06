package fp

import (
	"fmt"
	"iter"
	"slices"
)

func FindIter[T any](f func(T) bool) func(s iter.Seq[T]) Either[T] {
	return func(s iter.Seq[T]) Either[T] {
		for v := range s {
			if f(v) {
				return Right(v)
			}
		}
		return Left[T](fmt.Errorf("could not an find item"))
	}
}

func FindIter2[T1, T2 any](f func(T1, T2) bool) func(s iter.Seq2[T1, T2]) Either[T2] {
	return func(s iter.Seq2[T1, T2]) Either[T2] {
		for v1, v2 := range s {
			if f(v1, v2) {
				return Right(v2)
			}
		}
		return Left[T2](fmt.Errorf("could not find an item"))
	}
}

func Find[T any](f func(int, T) bool) func(s []T) Either[T] {
	return func(s []T) Either[T] {
		return FindIter2(f)(slices.All(s))
	}
}
