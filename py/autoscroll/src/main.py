#!/usr/bin/env python3

from sys import exit as sys_exit

from .autoscroll import Autoscroll


def start() -> None:
    sys_exit(Autoscroll().start(parse_argv=True))


if __name__ == "__main__":
    start()
