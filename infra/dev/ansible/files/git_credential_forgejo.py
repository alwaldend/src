#!/usr/bin/python3
"""Supply the managed user's fj login to Git for the configured Forgejo origin."""

import json
import os
import pathlib
import stat
import subprocess
import sys


def _run() -> None:
    """Answer Git's get operation for the expected HTTPS repository only."""
    if sys.argv[1:] != ["get"]:
        return
    request: dict[str, str] = {}
    for line in sys.stdin:
        key, separator, value = line.rstrip("\n").partition("=")
        if not separator:
            break
        request[key] = value
    host: str = "forgejo.alwaldend.com"
    if (
        request.get("protocol") != "https"
        or request.get("host") != host
        or request.get("path") != "alwaldend/src.git"
    ):
        return

    try:
        subprocess.run(
            [pathlib.Path(__file__).with_name("fj"), "--host", f"https://{host}", "whoami"],
            check=True,
            stdin=subprocess.DEVNULL,
            stdout=subprocess.DEVNULL,
            stderr=subprocess.DEVNULL,
            timeout=15,
        )
        data_home: pathlib.Path = pathlib.Path(
            os.environ.get("XDG_DATA_HOME") or pathlib.Path.home() / ".local/share"
        )
        credentials: pathlib.Path = data_home / "forgejo-cli/keys.json"
        descriptor: int = os.open(
            credentials, os.O_RDONLY | os.O_NOFOLLOW | os.O_NONBLOCK
        )
        with os.fdopen(descriptor, encoding="utf-8") as stream:
            info: os.stat_result = os.fstat(stream.fileno())
            if (
                not stat.S_ISREG(info.st_mode)
                or info.st_uid != os.getuid()
                or stat.S_IMODE(info.st_mode) & 0o077
            ):
                return
            token: object = json.load(stream)["hosts"][host]["token"]
    except (OSError, subprocess.SubprocessError, KeyError, TypeError, ValueError):
        return

    if isinstance(token, str) and token and not any(c in token for c in "\r\n\0"):
        sys.stdout.write(f"username=simeonwarrenbot\npassword={token}\n\n")


if __name__ == "__main__":
    _run()
