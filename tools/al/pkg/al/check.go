package al

import "errors"

// Panic on errors
func Check(errs... error) {
    err := errors.Join(errs...)
    if err != nil {
        panic(err)
    }
}
