from os import environ as os_environ
import typing

from . import arguments, functions

SCROLLING_SPEED: int = 300
SCROLLING_ACCELERATION_DISTANCE: int = 10
SCROLLING_SLEEP_INTERVAL_INITIAL: float = 0.1
SCROLLING_DEAD_AREA: int = 50

BUTTONS_START: int = 2
BUTTONS_HOLD: bool = False

ICON_ENABLE: bool = False
ICON_SIZE: int = 30
ICON_PATH: str = "resources/img/icon.svg"
ICON_ERROR: str = (
    "icon is enabled (it is disabled by default), but the "
    "'pyside6' package is not installed. "
    "remove '--icon-enable' or install the package\n"
    "\npip install pyside6\n"
)

CONFIG_PATH: str = f'{os_environ.get("HOME")}/.config/autoscroll/config.txt'
CONFIG_ENABLE: bool = False
CONFIG_INTERVAL: int = 5
CONFIG_ERROR_ENABLE: str = (
    "you are trying to enable the config ('enable' is set to 'True'), but the path is not valid"
)
CONFIG_ERROR_PARSE: str = (
    "you are trying to parse the config file, but 'enable' is 'False'"
)


DEBUG_SCROLL: bool = False
DEBUG_CLICK: bool = False
DEBUG_INITIAL: bool = False
DEBUG_FILE: bool = False
DEBUG_PADDING: int = 16

COORDINATE_NAME: str = "coordinates"


PARSER_INITIALIZER: dict[str, typing.Any] = {
    "prog": "linux-xorg-autoscroll",
    "formatter_class": arguments.ArgparseFormatter,
    "description": functions.get_resource_content("resources/txt/prolog.txt"),
    "fromfile_prefix_chars": "@",
}

ARGUMENTS: dict[str, typing.Any] = {
    "scrolling": {
        "speed": {
            "type": int,
            "help": (
                "R|constant part of the scrolling speed\n"
                f"[default: {SCROLLING_SPEED}]"
            ),
        },
        "dead-area": {
            "type": int,
            "help": (
                "R|size of the square area aroung the starting "
                "point where scrolling will stop, in pixels\n"
                f"[default: {SCROLLING_DEAD_AREA}]"
            ),
        },
        "acceleration": {
            "type": int,
            "help": (
                "R|dynamic part of the  scrolling speed, depends on the "
                "distance from the point where the scrolling started, "
                "can be set to 0\n"
                f"[default: {SCROLLING_ACCELERATION_DISTANCE}]"
            ),
        },
    },
    "buttons": {
        "hold": {
            "action": "store_const",
            "const": True,
            "help": (
                "R|if set, the scrolling will end once you release "
                "--buttons-start"
            ),
        },
        "start": {
            "type": int,
            "help": f"R|button that starts the scrolling\n[default: {BUTTONS_START}]",
        },
        "end": {
            "type": int,
            "help": (
                "R|button that ends the scrolling\n"
                "[default: --buttons-start]"
            ),
        },
    },
    "config": {
        "enable": {
            "action": "store_const",
            "const": True,
            "help": (
                "R|if set, arguments from the configuration file "
                "on --config-path will be loaded every --config-interval"
            ),
        },
        "path": {
            "type": str,
            "help": f"R|path to the configuration file\n[default: {CONFIG_PATH}]",
        },
        "interval": {
            "type": int,
            "help": (
                "R|how often the config file should be checked for "
                f"changes, in seconds\n[default: {CONFIG_INTERVAL}]"
            ),
        },
    },
    "icon": {
        "enable": {
            "action": "store_const",
            "const": True,
            "help": "if set, the icon will be enabled",
        },
        "path": {
            "type": str,
            "help": f"R|path to the icon\n[default: {ICON_PATH}]",
        },
        "size": {
            "type": int,
            "help": f"R|size of the icon, in pixels\n[default: {ICON_SIZE}]",
        },
    },
    "debug": {
        "file": {
            "action": "store_const",
            "const": True,
            "help": (
                "if set, every time the config file is parsed, "
                "information will be printed to stdout"
            ),
        },
        "click": {
            "action": "store_const",
            "const": True,
            "help": "if set, click info will be printed to stdout",
        },
        "scroll": {
            "action": "store_const",
            "const": True,
            "help": "if set, scroll info will be printed to stdout",
        },
        "initial": {
            "action": "store_const",
            "const": True,
            "help": "if set, startup configuration will be printed to stdout",
        },
    },
}
