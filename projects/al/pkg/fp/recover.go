package fp

import (
	"fmt"
	"runtime/debug"
)

func Recover[T1, T2 any, R Result[T2]](f func(T1) R) func(T1) Either[T2] {
	return func(v T1) (res Either[T2]) {
		defer func() {
			if r := recover(); r != nil {
				res = Left[T2](fmt.Errorf("panic: %s\n%s", r, string(debug.Stack())))
			}
		}()
		return EitherOf(Get(f(v)))
	}
}
