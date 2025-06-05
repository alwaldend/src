---
title: Shell
---

## Deactivate venv in subshell

A subshell doesn't have the `deactivate` function, so you have to use a
workaround:

```sh
. .venv/bin/activate
(
    echo "subshell, still in venv: $(which python)" && \
    . "$(dirname "$(which python)")/activate" && \
    deactivate && \
    echo "subshell, out of venv: $(which python)"
)
```

## A subshell `()` is different from a subprocess

https://unix.stackexchange.com/a/138498

A subshell starts out as an almost identical copy of the original shell process.
Under the hood, the shell calls the fork system call1, which creates a new
process whose code and memory are copies2. When the subshell is created, there
are very few differences between it and its parent. In particular, they have the
same variables. Even the $$ special variable keeps the same value in subshells:
it's the original shell's process ID. Similarly $PPID is the PID of the parent
of the original shell.

A few shells change a few variables in the subshell. Bash ≥4.0 sets BASHPID to
the PID of the shell process, which changes in subshells. Bash, zsh and mksh
arrange for $RANDOM to yield different values in the parent and in the subshell.
But apart from built-in special cases like these, all variables have the same
value in the subshell as in the original shell, the same export status, the same
read-only status, etc. All function definitions, alias definitions, shell
options and other settings are inherited as well.

A subshell created by (…) has the same file descriptors as its creator. Some
other means of creating subshells modify some file descriptors before executing
user code; for example, the left-hand side of a pipe runs in a subshell3 with
standard output connected to the pipe. The subshell also starts out with the
same current directory, the same signal mask, etc. One of the few exceptions is
that subshells do not inherit custom traps: ignored signals (trap '' SIGNAL)
remain ignored in the subshell, but other traps (trap CODE SIGNAL) are reset to
the default action4.

A subshell is thus different from executing a script. A script is a separate
program. This separate program might coincidentally be also a script which is
executed by the same interpreter as the parent, but this coincidence doesn't
give the separate program any special visibility on internal data of the parent.
Non-exported variables are internal data, so when the interpreter for the child
shell script is executed, it doesn't see these variables. Exported variables,
i.e. environment variables, are transmitted to executed programs.
