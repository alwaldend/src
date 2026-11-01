import sys

from autoscroll import internal


def cli() -> None:
    sys.exit(internal.Autoscroll().start(parse_argv=True))
