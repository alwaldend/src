import sys

from . import _internal


def cli() -> None:
    sys.exit(_internal.Autoscroll().start(parse_argv=True))
