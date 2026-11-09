package script

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	lua "github.com/yuin/gopher-lua"
	"github.com/yuin/gopher-lua/parse"
)

// Failure records a script stage while retaining its original cause locally.
func Failure(stage, filename string, cause error) error {
	return &evaluationError{stage: stage, filename: filename, cause: cause}
}

// evaluationError retains the original failure while classifying safe diagnostics.
type evaluationError struct {
	// stage distinguishes syntax, top-level execution and selection evaluation.
	stage string
	// filename identifies trusted interpreter stack frames; diagnostics omit its value.
	filename string
	// cause retains original error identity for local callers.
	cause error
}

// Error preserves local diagnostic context without changing error identity.
func (e *evaluationError) Error() string { return e.cause.Error() }

// Unwrap allows callers to inspect the original interpreter failure.
func (e *evaluationError) Unwrap() error { return e.cause }

// Diagnostic returns only the failure stage and source location, never error contents.
func Diagnostic(err error) string {
	var failure *evaluationError
	if !errors.As(err, &failure) {
		return "script evaluation failed"
	}
	message := "script " + failure.stage + " error"
	var api *lua.ApiError
	if !errors.As(err, &api) {
		return message
	}
	var syntax *parse.Error
	if errors.As(api.Cause, &syntax) {
		if syntax.Pos.Line > 0 {
			return fmt.Sprintf("%s at line %d", message, syntax.Pos.Line)
		}
		return message + " at end of input"
	}
	for _, frame := range strings.Split(api.StackTrace, "\n") {
		frame = strings.TrimSpace(frame)
		if tail, ok := strings.CutPrefix(frame, failure.filename+":"); ok {
			number, _, _ := strings.Cut(tail, ":")
			if line, err := strconv.Atoi(number); err == nil && line > 0 {
				return fmt.Sprintf("%s at line %d", message, line)
			}
		}
	}
	return message
}
