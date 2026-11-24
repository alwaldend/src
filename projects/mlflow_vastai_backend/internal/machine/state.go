package machine

import (
	"context"
	"errors"
	"fmt"

	"github.com/qmuntal/stateless"
)

// State identifies a state in an application machine.
type State string

// Event identifies a trigger in an application machine.
type Event string

// scope carries the state and typed inputs of a single invocation or transaction.
type scope[T any] struct {
	// state is the slot that Stateless reads and updates.
	state *State
	// input holds the values needed by the typed action.
	input T
}

// Machine binds a Stateless definition to the type of input its actions receive.
type Machine[T any] struct {
	// StateMachine configures and runs transitions through the scoped state.
	*stateless.StateMachine
}

// New creates a machine whose invocations supply typed state and action inputs.
func New[T any](initial State) *Machine[T] {
	machine := &Machine[T]{}
	// A serving invocation can block while separate job transactions fire.
	machine.StateMachine = stateless.NewStateMachineWithExternalStorage(func(ctx context.Context) (stateless.State, error) {
		if current, ok := ctx.Value(machine).(scope[T]); ok {
			return *current.state, nil
		}
		// ToGraph queries the initial state without an invocation context.
		return initial, nil
	}, func(ctx context.Context, next stateless.State) error {
		current, ok := ctx.Value(machine).(scope[T])
		if !ok {
			return errors.New("state machine event requires an invocation scope")
		}
		state, ok := next.(State)
		if !ok {
			return fmt.Errorf("invalid application state type %T", next)
		}
		*current.state = state
		return nil
	}, stateless.FiringImmediate)
	return machine
}

// WithInvocation binds the state and typed input used by subsequent transitions.
func (machine *Machine[T]) WithInvocation(ctx context.Context, state *State, input T) context.Context {
	return context.WithValue(ctx, machine, scope[T]{state: state, input: input})
}

// Input returns the typed input of the current invocation.
func (machine *Machine[T]) Input(ctx context.Context) (T, error) {
	current, ok := ctx.Value(machine).(scope[T])
	if !ok {
		return current.input, errors.New("state machine action requires an invocation scope")
	}
	return current.input, nil
}

// RunAction passes the current invocation's typed input to an action.
func (machine *Machine[T]) RunAction(ctx context.Context, action func(context.Context, T) error) error {
	input, err := machine.Input(ctx)
	if err != nil {
		return fmt.Errorf("read action input: %w", err)
	}
	if err := action(ctx, input); err != nil {
		return fmt.Errorf("run action: %w", err)
	}
	return nil
}
