"""Observe a wrapped Bazel process, runfiles, environment, and resource limits."""

import argparse
import json
import os
import pathlib

from python.runfiles import runfiles


def _main() -> None:
    """Write repeatable evidence without modifying host configuration."""
    parser: argparse.ArgumentParser = argparse.ArgumentParser(__doc__)
    parser.add_argument("--fixture", required=True)
    parser.add_argument("--report", required=True)
    parser.add_argument("--exit-code", type=int, default=0)
    arguments: argparse.Namespace
    remaining: list[str]
    arguments, remaining = parser.parse_known_args()
    resolver: runfiles.Runfiles | None = runfiles.Create()
    if resolver is None:
        raise RuntimeError("Runfiles are unavailable")
    fixture: str | None = resolver.Rlocation(arguments.fixture)
    if fixture is None:
        raise FileNotFoundError(arguments.fixture)
    relative: str = next(
        line[3:]
        for line in pathlib.Path("/proc/self/cgroup").read_text().splitlines()
        if line.startswith("0::")
    )
    scope: pathlib.Path = pathlib.Path("/sys/fs/cgroup") / relative.lstrip("/")
    report: dict[str, object] = {
        "arguments": remaining,
        "fixture": pathlib.Path(fixture).read_text(),
        "environment": os.environ.get("SYSTEMD_RUN_TEST_VALUE"),
        "cwd": os.getcwd(),
        "pid": os.getpid(),
        "limits": {
            name: (scope / name).read_text().strip()
            for name in ["memory.high", "memory.max", "memory.swap.max"]
        },
    }
    pathlib.Path(arguments.report).write_text(
        json.dumps(report, indent=2) + "\n"
    )
    raise SystemExit(arguments.exit_code)


if __name__ == "__main__":
    _main()
