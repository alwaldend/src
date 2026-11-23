package workload

import (
	"context"
	"errors"
	"fmt"
	"os/exec"
	"syscall"

	"golang.org/x/sys/unix"
)

// waitForCommand observes a child exit while keeping its PID reserved until Wait.
func waitForCommand(pid int, flags int) (bool, error) {
	var info unix.Siginfo
	for {
		err := unix.Waitid(unix.P_PID, pid, &info, unix.WEXITED|unix.WNOWAIT|flags, nil)
		if errors.Is(err, unix.EINTR) {
			continue
		}
		if err != nil {
			return false, fmt.Errorf("observe process %d: %w", pid, err)
		}
		return info.Signo != 0, nil
	}
}

// runCommand owns a command group until the command exit has been confirmed.
func runCommand(ctx context.Context, command *exec.Cmd) (Outcome, error) {
	if err := ctx.Err(); err != nil {
		return Outcome{}, fmt.Errorf("before process start: %w", err)
	}
	command.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	if err := command.Start(); err != nil {
		return Outcome{}, fmt.Errorf("start process: %w", err)
	}
	outcome := Outcome{Started: true}
	pid := command.Process.Pid
	exited := make(chan error, 1)
	go func() {
		_, err := waitForCommand(pid, 0)
		exited <- err
	}()
	var observeErr, signalErr error
	var signalled bool
	select {
	case observeErr = <-exited:
	case <-ctx.Done():
		// WNOWAIT leaves the PID occupied, so a group signal cannot target a
		// new process after the child exits concurrently with cancellation.
		alreadyExited, err := waitForCommand(pid, unix.WNOHANG)
		if err != nil {
			signalErr = err
		} else if !alreadyExited {
			err := unix.Kill(-pid, unix.SIGKILL)
			signalled = err == nil
			if err != nil && !errors.Is(err, unix.ESRCH) {
				signalErr = fmt.Errorf("stop process group %d: %w", pid, err)
			}
		}
		observeErr = <-exited
	}
	waitErr := command.Wait()
	if command.ProcessState != nil {
		outcome.Observed = true
		if code := command.ProcessState.ExitCode(); code >= 0 {
			value := int32(code)
			outcome.ExitCode = &value
		}
		status, ok := command.ProcessState.Sys().(syscall.WaitStatus)
		outcome.Stopped = ok && signalled && status.Signaled() && status.Signal() == unix.SIGKILL
	}
	if waitErr != nil {
		waitErr = fmt.Errorf("wait for process %d: %w", pid, waitErr)
	}
	return outcome, errors.Join(observeErr, signalErr, waitErr)
}
