package fp

import (
	"fmt"
	"io"
	"os"
)

type Result[T any] interface {
	Get() (T, error)
}

type ResultFunc[R any] = func(r Result[R]) (R, error)
type ResultFuncSuccess[R any] = func(r R) (R, error)
type ResultFuncFail[R any] = func(r R, err error) (R, error)

func Get[R any](r Result[R]) (R, error) {
	return r.Get()
}

func OnError[T1, T2 any, R Result[T2]](f func(T1) R, fe func(err error) error) func(T1) Either[T2] {
	return func(v T1) Either[T2] {
		res, err := Get(f(v))
		if err != nil {
			err = fe(err)
			if err != nil {
				return Left[T2](err)
			}
		}
		return Right(res)
	}
}

func Errorf(s string, a ...any) func(error) error {
	return func(err error) error {
		a = append(a, err)
		return fmt.Errorf(s, a...)
	}
}

func Fold[R any](f1 ResultFuncSuccess[R], f2 ResultFuncFail[R]) ResultFunc[R] {
	return func(r Result[R]) (R, error) {
		val, err := Get(r)
		if err != nil {
			return f2(val, err)
		}
		return f1(val)
	}
}

func GetOrExit[R any](out io.Writer, code int) ResultFunc[R] {
	return Fold(
		ReturnE,
		func(r R, err error) (R, error) {
			fmt.Fprintf(out, "%s\n", err)
			os.Exit(code)
			return r, fmt.Errorf("did not exit for some reason: %w", err)
		},
	)
}

func GetOrPanic[R any](r Result[R]) {
	Fold(ReturnE, func(r R, err error) (R, error) { panic(err) })(r)
}
