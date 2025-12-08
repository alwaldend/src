import os
import threading
import time
import typing

import pynput.mouse

from .arguments import ArgparseParser, parse_arguments
from .constants import (
    ARGUMENTS,
    BUTTONS_HOLD,
    BUTTONS_START,
    CONFIG_ENABLE,
    CONFIG_ERROR_ENABLE,
    CONFIG_INTERVAL,
    CONFIG_PATH,
    COORDINATE_NAME,
    DEBUG_CLICK,
    DEBUG_FILE,
    DEBUG_INITIAL,
    DEBUG_PADDING,
    DEBUG_SCROLL,
    ICON_ENABLE,
    ICON_ERROR,
    ICON_PATH,
    ICON_SIZE,
    PARSER_INITIALIZER,
    SCROLLING_ACCELERATION_DISTANCE,
    SCROLLING_DEAD_AREA,
    SCROLLING_SLEEP_INTERVAL_INITIAL,
    SCROLLING_SPEED,
)
from .functions import (
    construct_coordinates,
    convert_bool,
    has_dict,
    raise_type_error,
    return_kwargs,
    return_none,
)


class Base:

    debug_keys_ignore = []
    debug_keys_only = []

    def update(self, *args, **kwargs) -> None: ...

    def json(self) -> dict[str, typing.Any]: ...

    def __init__(self, *args, **kwargs) -> None:
        self.update(*args, **kwargs)

    def _set_if_nonexistent(self, name: str, value: typing.Any) -> None:
        return None if hasattr(self, name) else setattr(self, name, value)

    @staticmethod
    def _convert_callable(value: typing.Any) -> typing.Any:
        return value

    def _set(
        self,
        _name: str,
        _nonexistent_value: typing.Any,
        _value: typing.Any,
        _type: typing.Any = type(None),
        _callable: typing.Callable = _convert_callable,
        **kwargs,
    ) -> None:
        self._set_if_nonexistent(_name, _nonexistent_value)
        value = self._convert(
            _value, getattr(self, _name), _type, _callable, **kwargs
        )
        setattr(self, _name, value)

    def _convert(
        self,
        _value: typing.Any,
        none_value: typing.Any = None,
        _type: typing.Any = type(None),
        _callable: typing.Callable = _convert_callable,
        **kwargs,
    ) -> typing.Any:
        """
        converts

        _value is _type and _callable is callable -> _callable(_value, **kwargs)

        _value is None -> none_value

        _value is _type -> _callable is not callable error

        _value is not None or _type type error
        """
        if isinstance(_value, _type) and callable(_callable):
            return _callable(_value, **kwargs)
        if _value is None:
            return none_value
        if isinstance(_value, _type):
            raise ValueError(
                f'value "{_value}" is "{_type}", but the '
                "conversion function is not callable"
            )
        return raise_type_error(_value, (type(None), _type))

    def _loop(
        self,
        condition: typing.Callable = return_none,
        action: typing.Callable = return_none,
        condition_getter: typing.Callable = return_kwargs,
        action_getter: typing.Callable = return_kwargs,
        condition_parameters: dict[str, typing.Any] = {},
        action_parameters: dict[str, typing.Any] = {},
    ) -> None:

        if False in list(
            map(callable, (condition, action, condition_getter, action_getter))
        ):
            raise TypeError(f"some functions are not callable")
        if not has_dict(condition_parameters, action_parameters):
            raise TypeError("parameters should have '__dict__' attribute")
        while condition(**condition_getter(**condition_parameters)):
            action(**action_getter(**action_parameters))

    def _print(
        self,
        header: str,
        do_print: bool = True,
        keys_only: typing.Optional[list[str]] = None,
        keys_ignore: typing.Optional[ list[str]] = None,
    ) -> str:
        result = f"\n[{header}]\n{self._debug(keys_only, keys_ignore)}"
        if do_print:
            print(result)
        return result

    def __repr__(self) -> str:
        return self._debug()

    def __str__(self) -> str:
        return self._debug()

    def _debug(
        self, keys_only: typing.Optional[list[str]] = None, keys_ignore: typing.Optional[list[str]] = None
    ) -> str:
        name = getattr(self, "name") if hasattr(self, "name") else type(self).__name__
        raise_type_error(name, str)
        debug = []
        not_base = {}
        for key, value in self.json().items():
            if not self._debug_key_is_valid(key, keys_only, keys_ignore):
                continue
            if isinstance(value, Base):
                debug.append(str(value))
                continue
            not_base[key] = value if isinstance(value, str) else str(value)
        if not_base:
            debug.append(self._construct_debug(name.capitalize(), **not_base))
        return "\n".join(debug)

    def _debug_key_is_valid(
        self,
        name: str,
        keys_only: typing.Optional[list[str]] = None,
        keys_ignore: typing.Optional[list[str]] = None,
    ) -> bool:
        only = self.debug_keys_only if keys_only is None else keys_only
        ignore = self.debug_keys_ignore if keys_ignore is None else keys_ignore
        return name in only if only else name not in ignore

    def _construct_debug(self, _header: str, *args, **kwargs) -> str:
        result = ", ".join(
            [item for item in args]
            + [f"{name} - {value}" for name, value in kwargs.items()]
        )
        return _header.ljust(DEBUG_PADDING, " ") + result


class Coordinate(Base):

    def update(
        self,
        current: typing.Optional[typing.Union[int, str]] = None,
        previous: typing.Optional[ typing.Union[int, str]] = None,
        initial: typing.Optional[ typing.Union[int, str]] = None,
    ) -> None:
        self.initial = initial
        self.previous = previous
        self.current = current

    def distance(self, absolute: bool = False) -> int:
        distance = self.initial - self.current
        return abs(distance) if absolute else distance

    def direction(self) -> int:
        # direction > 0 -> 1, direction == 0 -> 0, direction < 0 -> -1
        return (self.distance() == 0) + (1 if self.distance() > 0 else -1)

    @property
    def current(self) -> int:
        return getattr(self, "_current")


    @property
    def previous(self) -> int:
        return getattr(self, "_previous")

    @property
    def initial(self) -> int:
        return getattr(self, "_initial")

    @current.setter
    def current(self, value: typing.Optional[typing.Union[int, str]] = None) -> None:
        _current = getattr(self, "_current", None)
        self._set("_current", 0, value, (str, int), int)
        self._set("_previous", 0, _current, (str, int), int)

    @previous.setter
    def previous(self, value: typing.Optional[typing.Union[int, str]] = None) -> None:
        self._set("_previous", 0, value, (str, int), int)

    @initial.setter
    def initial(self, value: typing.Optional[typing.Union[int, str]] = None) -> None:
        self._set("_initial", 0, value, (str, int), int)

    def json(self) -> dict[str, typing.Any]:
        return {
            "current": self.current,
            "previous": self.previous,
            "initial": self.initial,
        }


class Coordinates(Base):

    def update(
        self,
        x: typing.Optional[typing.Union[Coordinate, str, int]] = None,
        y: typing.Optional[ typing.Union[Coordinate, str, int]] = None,
        name: typing.Optional[str]= None,
    ) -> None:
        self.x = x
        self.y = y
        self.name: str = self._convert(name, COORDINATE_NAME, str)

    def direction(self) -> tuple[int, int]:
        return (self.x.direction(), self.y.direction())

    def distance(
        self, absolute: bool = False, reverse: bool = False
    ) -> tuple[int, int]:
        x, y = self.x.distance(absolute), self.y.distance(absolute)
        return (-x, -y) if reverse else (x, y)

    @property
    def x(self) -> Coordinate:
        return getattr(self, "_x")

    @property
    def y(self) -> Coordinate:
        return getattr(self, "_y")

    @property
    def current(self) -> tuple[int, int]:
        return (self.x.current, self.y.current)

    @property
    def initial(self) -> tuple[int, int]:
        return (self.x.initial, self.y.initial)

    @property
    def previous(self) -> tuple[int, int]:
        return (self.x.previous, self.y.previous)

    @y.setter
    def y(self, value: typing.Union[str, int, Coordinate, None] = None) -> None:
        self._set(
            "_y",
            Coordinate(),
            value,
            (Coordinate, str, int),
            self._convert_coordinate,
            name="_y",
        )

    @x.setter
    def x(self, value: typing.Union[str, int, Coordinate, None] = None) -> None:
        self._set(
            "_x",
            Coordinate(),
            value,
            (Coordinate, str, int),
            self._convert_coordinate,
            name="_x",
        )

    @current.setter
    def current(self, value: tuple[typing.Union[str, int], typing.Union[str, int]]) -> None:
        self.x.current, self.y.current = self._convert(
            value, (None, None), typing.Iterable, self._convert_iterable
        )

    @initial.setter
    def initial(self, value: tuple[typing.Union[str, int], typing.Union[str, int]]) -> None:
        self.x.initial, self.y.initial = self._convert(
            value, (None, None), typing.Iterable, self._convert_iterable
        )

    @previous.setter
    def previous(self, value: tuple[typing.Union[str, int], typing.Union[str, int]]) -> None:
        self.x.previous, self.y.previous = self._convert(
            value, (None, None), typing.Iterable, self._convert_iterable
        )

    def _convert_coordinate(
        self, value: typing.Union[Coordinate, int], name: str
    ) -> Coordinate:
        if isinstance(value, (int, str)):
            coordinate = getattr(self, name)
            setattr(coordinate, "current", int(value))
            return coordinate
        return raise_type_error(value, Coordinate)

    def _convert_iterable(self, value: typing.Any) -> tuple[typing.Optional[int], typing.Optional[int]]:
        raise_type_error(value, typing.Iterable)
        x, y = None, None
        if len(value) >= 1:
            x = self._convert(value[0], None, (str, int), int)
        if len(value) >= 2:
            y = self._convert(value[1], None, (str, int), int)
        return x, y

    def json(self) -> dict[str, typing.Any]:
        coordinates = {
            "direction": str(self.direction()),
            "distance": str(self.distance()),
        }
        coordinates.update(
            {
                item: construct_coordinates(
                    getattr(self.x, item), getattr(self.y, item)
                )
                for item in ("current", "previous", "initial")
            }
        )
        return coordinates


class Buttons(Base):

    def __init__(self, *args, **kwargs) -> None:
        self.button: typing.Optional[pynput.mouse.Button] = None
        self.is_pressed: typing.Optional[bool] = None
        self.update(*args, **kwargs)

    def update(
        self,
        start: typing.Optional[typing.Union[pynput.mouse.Button, int, str]] = None,
        end: typing.Optional[ typing.Union[pynput.mouse.Button, int, str]] = None,
        hold: typing.Optional[ typing.Union[bool, str]] = None,
    ) -> None:
        self.start = start
        self.end = end
        self.hold = hold

    def press(self, button: pynput.mouse.Button, pressed: bool) -> None:
        self.button, self.is_pressed = button, pressed

    def press_clear(self) -> None:
        self.button, self.is_pressed = None, None

    def is_start(self) -> bool:
        return self.start == self.button

    def is_end(self) -> bool:
        return self.end == self.button

    def was_action(self) -> bool:
        return self.button is not None and self.is_pressed is not None

    def was_start_pressed(self) -> bool:
        return self.was_action() and self.is_start() and bool(self.is_pressed)

    def was_end_pressed(self) -> bool:
        return self.was_action() and self.is_end() and bool(self.is_pressed)

    def was_start_released(self) -> bool:
        return self.was_action() and self.is_start() and not self.is_pressed

    def was_start_released_with_hold(self) -> bool:
        return self.was_start_released() and self.hold

    def was_end_released(self) -> bool:
        return self.was_action() and self.is_end() and not self.is_pressed

    @property
    def start(self) -> pynput.mouse.Button:
        return getattr(self, "_start")

    @property
    def hold(self) -> bool:
        return getattr(self, "_hold")

    @property
    def end(self) -> pynput.mouse.Button:
        return getattr(self, "_end")

    @start.setter
    def start(self, value: typing.Optional[typing.Union[int, str, pynput.mouse.Button]] = None) -> None:
        self._set(
            "_start",
            pynput.mouse.Button(BUTTONS_START),
            value,
            (int, pynput.mouse.Button, str),
            self._convert_button,
        )

    @end.setter
    def end(self, value: typing.Optional[typing.Union[int, str, pynput.mouse.Button]] = None) -> None:
        self._set(
            "_end", self.start, value, (int, pynput.mouse.Button, str), self._convert_button
        )

    @hold.setter
    def hold(self, value: typing.Optional[typing.Union[str, bool]] = None) -> None:
        self._set("_hold", BUTTONS_HOLD, value, (str, bool), convert_bool)

    @staticmethod
    def _convert_button(value: typing.Union[pynput.mouse.Button, int, str]) -> pynput.mouse.Button:
        return value if isinstance(value, pynput.mouse.Button) else pynput.mouse.Button(int(value))

    def json(self) -> dict[str, typing.Any]:
        return {
            "start": self.start.name,
            "end": self.end.name,
            "hold": self.hold,
            "pressed button": self.button,
            "pressed": self.is_pressed,
        }


class Scrolling(Base):

    def __init__(self, *args, **kwargs) -> None:
        self.sleep_interval: float = SCROLLING_SLEEP_INTERVAL_INITIAL
        self.controller: pynput.mouse.Controller = pynput.mouse.Controller()
        self.coordinates: Coordinates = Coordinates()
        self.coordinates.debug_keys_ignore = ["direction"]
        self.direction: Coordinates = Coordinates(name="direction")
        self.direction.debug_keys_only = ["direction"]

        self.event_end: threading.Event = threading.Event()
        self.event_scrolling: threading.Event = threading.Event()
        self.event_started: threading.Event = threading.Event()
        self.event_ended: threading.Event = threading.Event()

        self.update(*args, **kwargs)

    def update(
        self,
        dead_area: typing.Optional[typing.Union[str, int]] = None,
        speed: typing.Optional[typing.Union[str, int]] = None,
        acceleration: typing.Optional[typing.Union[str, int]] = None,
    ) -> None:
        self.speed = speed
        self.dead_area = dead_area
        self.acceleration = acceleration

    def sleep_for_interval(self) -> None:
        time.sleep(self.sleep_interval)

    def wait(self) -> None:
        self.event_scrolling.wait()

    def scroll_once(self) -> None:
        self.controller.scroll(*self.direction.current)

    def start(self) -> None:
        self.event_started.set()
        self.event_scrolling.set()

    def stop(self) -> None:
        self.event_scrolling.clear()
        self.event_ended.set()

    def clear_started_and_ended(self) -> None:
        self.event_ended.clear()
        self.event_started.clear()

    def is_scrolling(self) -> bool:
        return self.event_scrolling.is_set()

    def has_started(self) -> bool:
        return self.event_started.is_set()

    def has_ended(self) -> bool:
        return self.event_ended.is_set()

    def is_not_end(self) -> bool:
        return not self.event_end.is_set()

    def is_dead_area(self) -> bool:
        distance = self.coordinates.distance(absolute=True)
        return distance[0] <= self.dead_area and distance[1] <= self.dead_area

    def set_interval(self) -> None:
        distance = self.coordinates.distance(absolute=True)
        interval = self.acceleration * max(distance) + self.speed
        self.sleep_interval = (
            abs(100 / interval)
            if interval
            else SCROLLING_SLEEP_INTERVAL_INITIAL
        )

    def set_initial_coordinates(self, x: int, y: int) -> None:
        self.coordinates.initial = x, y

    def set_direction_and_coordinates(self, x: int, y: int) -> None:
        self.coordinates.update(x, y)
        self.direction.current = (
            (0, 0) if self.is_dead_area() else self.coordinates.direction()
        )

    def json(self) -> dict[str, typing.Any]:
        return {
            "active": self.is_scrolling(),
            "interval": self.sleep_interval,
            "acceleration": self.acceleration,
            "dead_area": self.dead_area,
            "started": self.has_started(),
            "ended": self.has_ended(),
            "coordinates": self.coordinates,
            "direction": self.direction,
        }

    @property
    def speed(self) -> int:
        return getattr(self, "_speed")

    @property
    def dead_area(self) -> int:
        return getattr(self, "_dead_area")

    @property
    def acceleration(self) -> int:
        return getattr(self, "_acceleration")

    @speed.setter
    def speed(self, value: typing.Optional[typing.Union[str, int]] = None) -> None:
        self._set("_speed", SCROLLING_SPEED, value, (str, int), int)

    @dead_area.setter
    def dead_area(self, value: typing.Optional[typing.Union[str, int]] = None) -> None:
        self._set("_dead_area", SCROLLING_DEAD_AREA, value, (str, int), int)

    @acceleration.setter
    def acceleration(self, value: typing.Optional[typing.Union[str, int]]) -> None:
        self._set(
            "_acceleration",
            SCROLLING_ACCELERATION_DISTANCE,
            value,
            (str, int),
            int,
        )


class Icon(Base):

    def __init__(self, *args, **kwargs) -> None:
        self.application: typing.Optional[typing.Any] = None
        self.event_icon_enabled: threading.Event = threading.Event()
        self.event_qt_application_started: threading.Event = threading.Event()
        self.update(*args, **kwargs)

    def update(
        self,
        enable: typing.Optional[typing.Union[str, bool]] = None,
        path: typing.Optional[str ]= None,
        size: typing.Optional[typing.Union[str, int]] = None,
    ) -> None:
        self.path = path
        self.size = size
        self.enable = enable
        self.icon = self.path, self.size

    def show(self, x: int, y: int) -> None:
        return self.icon.show(x, y) if self.enable and self.icon else None

    def hide(self) -> None:
        return self.icon.hide() if self.enable and self.icon else None

    def json(self) -> dict[str, typing.Any]:
        return {"enable": self.enable, "path": self.path, "size": self.size}

    def start_qt_when_icon_is_enabled(self) -> None:
        self.event_icon_enabled.wait()
        self.application = self._get_qt(True)
        self.event_qt_application_started.set()
        self.application.exec()

    @property
    def path(self) -> str:
        return getattr(self, "_path")

    @property
    def size(self) -> int:
        return getattr(self, "_size")

    @property
    def enable(self) -> bool:
        return getattr(self, "_enable")

    @property
    def icon(self) -> typing.Optional[typing.Any]:
        return getattr(self, "_icon")

    @path.setter
    def path(self, value: typing.Optional[str] = None) -> None:
        self._set("_path", ICON_PATH, value, str)

    @size.setter
    def size(self, value: typing.Optional[typing.Union[str, int]] =None) -> None:
        self._set("_size", ICON_SIZE, value, (str, int), int)

    @enable.setter
    def enable(self, value: typing.Optional[typing.Union[str, bool]] =None) -> None:
        self._set("_enable", ICON_ENABLE, value, (str, bool), convert_bool)
        if self.enable:
            self.event_icon_enabled.set()

    @icon.setter
    def icon(self, value: typing.Optional[tuple[str, int]] = None) -> None:
        if not self.enable:
            self._icon = None
            return
        self.event_qt_application_started.wait()
        if not isinstance(value, tuple) or not self.icon or not self.application:
            raise Exception("missing value, icon, or application: "
                f"{value}, {self.icon}, {self.application}")
        if getattr(self, "_icon", None) is not None:
            self.icon.update_icon(*value)
            return
        self._icon = self._get_qt()(*value)
        self.application.setActiveWindow(self.icon)

    def _get_qt(self, get_application: bool = False) -> typing.Callable:
        try:
            from .qt import Icon as qt_icon
            from .qt import application
        except ImportError as exception:
            raise ValueError(ICON_ERROR) from exception
        return application if get_application else qt_icon


class Debug(Base):

    def update(
        self,
        scroll: typing.Optional[bool] = None,
        file: typing.Optional[bool] = None,
        click: typing.Optional[bool] = None,
        initial: typing.Optional[bool] = None,
    ) -> None:
        self.scroll= scroll
        self.click = click
        self.initial = initial
        self.file = file

    def json(self) -> dict[str, typing.Any]:
        return {
            "scroll": self.scroll,
            "click": self.click,
            "initial": self.initial,
            "file": self.file,
        }

    @property
    def scroll(self) -> bool:
        return getattr(self, "_scroll")

    @property
    def click(self) -> bool:
        return getattr(self, "_click")

    @property
    def initial(self) -> bool:
        return getattr(self, "_initial")

    @property
    def file(self) -> bool:
        return getattr(self, "_file")

    @scroll.setter
    def scroll(self, value: typing.Optional[bool] = None) -> None:
        self._set("_scroll", DEBUG_SCROLL, value, (str, bool), convert_bool)

    @click.setter
    def click(self, value: typing.Optional[bool] = None) -> None:
        self._set("_click", DEBUG_CLICK, value, (str, bool), convert_bool)

    @initial.setter
    def initial(self, value: typing.Optional[bool] =None) -> None:
        self._set("_initial", DEBUG_INITIAL, value, (str, bool), convert_bool)

    @file.setter
    def file(self, value: typing.Optional[bool] =None) -> None:
        self._set("_file", DEBUG_FILE, value, (str, bool), convert_bool)


class Config(Base):

    debug_keys_ignore = "content"

    def __init__(self, *args, **kwargs) -> None:
        self._stamp: float = 0
        self.event_enabled: threading.Event = threading.Event()
        self._parse_config_file_content: dict[str, typing.Any] = {}
        self.argument_parser: ArgparseParser = ArgparseParser(
            **PARSER_INITIALIZER
        ).add_arguments(**ARGUMENTS)
        self.update(*args, **kwargs)

    def update(
        self,
        enable: typing.Union[bool, str, None] = None,
        path: typing.Optional[str] = None,
        interval: typing.Union[str, int, None] = None,
    ) -> None:
        self.path = path
        self.enable = enable
        self.interval = interval

    def wait(self) -> None:
        self.event_enabled.wait()

    def parse_argv(self) -> dict[str, typing.Any]:
        return self._parse()

    def parse_string(self, value: str) -> dict[str, typing.Any]:
        return self._parse(value.split())

    def parse_config_file(self) -> dict[str, typing.Any]:
        self._parse_config_file_content = self._parse_config_file()
        return self._parse_config_file_content

    def _parse_config_file(self) -> dict[str, typing.Any]:
        if not self._has_file_changed():
            return {}
        with open(self.path, "r") as config_file:
            config = config_file.read()
        result = self.parse_string(config.replace("\n", " "))
        return result

    def _parse(self, *args, **kwargs) -> dict[str, typing.Any]:
        return parse_arguments(
            **vars(self.argument_parser.parse_args(*args, **kwargs))
        )

    def _has_file_changed(self) -> bool:
        stamp = os.stat(self.path).st_mtime
        if stamp == self._stamp:
            return False
        self._stamp = stamp
        return True

    @property
    def enable(self) -> bool:
        return getattr(self, "_enable")

    @property
    def path(self) -> str:
        return getattr(self, "_path")

    @property
    def interval(self) -> int:
        return getattr(self, "_interval")

    @enable.setter
    def enable(self, value: typing.Union[bool, str, None] = None) -> None:
        self._set("_enable", CONFIG_ENABLE, value, (bool, str), convert_bool)
        if self.enable and not self.path:
            self.enable = False
            raise ValueError(f"{CONFIG_ERROR_ENABLE}, path - {self.path}")
        if self.enable and not self.event_enabled.is_set():
            self.event_enabled.set()

    @interval.setter
    def interval(self, value: typing.Union[str, int, None] = None) -> None:
        self._set("_interval", CONFIG_INTERVAL, value, (int, str), int)

    @path.setter
    def path(self, value: typing.Optional[str] = None) -> None:
        self._set("_path", CONFIG_PATH, value, str)

    def json(self) -> dict[str, typing.Any]:
        return {
            "path": self.path,
            "enable": self.enable,
            "interval": self.interval,
            "content": self._parse_config_file_content,
        }
