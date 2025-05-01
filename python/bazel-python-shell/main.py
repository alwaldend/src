import os
import shlex
import sys
import typing

def main(argv: typing.Sequence[str]) -> None:
    bazel_working_dir = os.environ.get("BAZEL_WORKING_DIRECTORY")
    shell_type = os.environ.get("BAZEL_PYTHON_SHELL_TYPE", "python")
    if bazel_working_dir:
        os.chdir(bazel_working_dir)
    if shell_type == "python":
        executable = sys.executable
        args = argv 
    elif shell_type == "bash":
        executable = "/bin/bash"
        args = (executable, "-c", shlex.join(argv[1:]))
    elif shell_type == "sh":
        executable = "/bin/sh"
        args = (executable, "-c", shlex.join(argv[1:]))
    else:
       raise Exception(f"invalid shell_type: {shell_type}")
    os.execv(executable, args)

if __name__ == "__main__":
    main(sys.argv)
