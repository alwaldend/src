"""Substitute systemctl only inside an isolated Forgejo role integration test."""

import json
import os
import pathlib
import shlex
import shutil
import signal
import subprocess
import sys
import time


def _running(root: pathlib.Path) -> bool:
    """Check the fixture process without querying the host service manager."""
    try:
        pid: int = int((root / "pid").read_text())
        status: str = pathlib.Path(f"/proc/{pid}/stat").read_text()
        return status.split(") ", 1)[1][0] != "Z"
    except (FileNotFoundError, ProcessLookupError):
        return False


def _event(root: pathlib.Path, action: str) -> None:
    """Retain non-secret lifecycle evidence for the parent test."""
    with (root / "events").open("a") as output:
        output.write(json.dumps({"action": action}) + "\n")


def _main() -> int:
    """Implement the bounded systemctl commands used by the real Ansible module."""
    root: pathlib.Path = pathlib.Path(os.environ["FORGEJO_TEST_ROOT"])
    arguments: list[str] = [
        arg for arg in sys.argv[1:] if not arg.startswith("-")
    ]
    action: str = arguments[0]
    if arguments[1:] and arguments[1] != "forgejo.service":
        raise ValueError("fixture refuses an unexpected service")
    active: bool = _running(root)
    if action == "show":
        print("LoadState=loaded")
        print("ActiveState=" + ("active" if active else "inactive"))
        print("SubState=" + ("running" if active else "dead"))
        print(
            "UnitFileState="
            + ("enabled" if (root / "enabled").exists() else "disabled")
        )
        return 0
    if action == "is-active":
        print("active" if active else "inactive")
        return 0 if active else 3
    if action == "is-enabled":
        print("enabled" if (root / "enabled").exists() else "disabled")
        return 0 if (root / "enabled").exists() else 1
    if action == "enable":
        (root / "enabled").touch()
        return 0
    if action in ("daemon-reload", "reset-failed"):
        return 0
    if action == "stop":
        _event(root, "stop")
        if active:
            os.killpg(int((root / "pid").read_text()), signal.SIGTERM)
            for _ in range(200):
                if not _running(root):
                    break
                time.sleep(0.05)
            else:
                raise TimeoutError("fixture service did not stop")
        if (root / "fail-backup").exists():
            (root / "fail-backup").unlink()
            shutil.rmtree(root / "backups")
        if (root / "fail-marker").exists():
            (root / "fail-marker").unlink()
            (root / "backups/pending.json").mkdir(mode=0o700)
        return 0
    if action == "start":
        _event(root, "start")
        unit: str = (root / "forgejo.service").read_text()
        command: str = next(
            line[10:]
            for line in unit.splitlines()
            if line.startswith("ExecStart=")
        )
        directory: str = next(
            line[17:]
            for line in unit.splitlines()
            if line.startswith("WorkingDirectory=")
        )
        with (root / "server.log").open("ab") as output:
            child: subprocess.Popen[bytes] = subprocess.Popen(
                shlex.split(command),
                cwd=directory,
                stdin=subprocess.DEVNULL,
                stdout=output,
                stderr=subprocess.STDOUT,
                start_new_session=True,
            )
        (root / "pid").write_text(str(child.pid))
        return 0
    raise ValueError("fixture refuses an unexpected operation")


if __name__ == "__main__":
    sys.exit(_main())
