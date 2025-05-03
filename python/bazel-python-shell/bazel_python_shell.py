import sys
import time
import os
import shlex
import typing

def main(argv: typing.Sequence[str]) -> None:
    shell_type = os.environ.get("BAZEL_PYTHON_SHELL_TYPE", "python")
    if shell_type == "python":
        executable = sys.executable
        args = argv[1:]
    elif shell_type == "bash":
        executable = "/bin/bash"
        args = ("-c", shlex.join(argv[1:]))
    elif shell_type == "sh":
        executable = "/bin/sh"
        args = ("-c", shlex.join(argv[1:]))
    else:
       raise Exception(f"invalid shell_type: {shell_type}")
    sys.stdout.flush()
    os.execv(executable, (executable, *args))

if __name__ == "__main__":
    main(sys.argv)
