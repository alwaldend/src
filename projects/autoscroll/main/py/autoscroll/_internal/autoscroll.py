import threading
import typing
import time

import pynput.mouse

from . import support


class Autoscroll(support.Base):

    def __init__(self, *args, **kwargs) -> None:
        self.scrolling: support.Scrolling = support.Scrolling()
        self.icon: support.Icon = support.Icon()
        self.config: support.Config = support.Config()
        self.buttons: support.Buttons = support.Buttons()
        self.debug: support.Debug = support.Debug()
        self.event_end: threading.Event = threading.Event()

        # update from initializer arguments
        self.update(*args, **kwargs)

        # threads
        # mouse actions
        self.thread_scroll_listener = pynput.mouse.Listener(
            on_move=self._on_move, on_click=self._on_click, daemon=True
        )
        # scroll
        self.thread_scroll_action = threading.Thread(
            target=self._loop,
            daemon=True,
            args=(self._is_not_end, self._scroll),
        )
        # update from the config file (if enabled)
        self.thread_config = threading.Thread(
            target=self._loop,
            daemon=True,
            args=(self._is_not_end, self._update_from_config_file),
        )

    def start(self, parse_argv: bool = False) -> None:
        # update from the command line
        # it is a thread because in order to create a qt icon widget,
        # a qt application has to be running in the main thread.
        # so, if the icon is enabled, the thread will have to wait untill the
        # qt application is running in the main thread
        threading.Thread(
            target=self.update,
            kwargs=(self.config.parse_argv() if parse_argv else {}),
        ).start()
        # start listening for mouse movements and clicks
        self.thread_scroll_listener.start()
        # start the scrolling loop
        self.thread_scroll_action.start()
        # update from the command line
        # start listening for changes in the config file
        self.thread_config.start()
        # debug
        self._print("initial", self.debug.initial)
        # wait untill the icon is enabled, then run a qt application
        self.icon.start_qt_when_icon_is_enabled()

    def _update_from_config_file(self) -> None:
        self.config.wait()
        self.update(**self.config.parse_config_file())
        self.config._print("config", self.debug.file, ["content"])
        time.sleep(self.config.interval)

    def _scroll(self) -> None:
        self.scrolling.wait()
        self.scrolling.sleep_for_interval()
        self.scrolling.scroll_once()

    def _on_move(self, x: int, y: int) -> None:
        self.scrolling.set_direction_and_coordinates(x, y)
        self.scrolling.set_interval()

    def _on_click(self, x: int, y: int, button: pynput.mouse.Button, pressed: bool) -> None:
        self.buttons.press(button, pressed)

        if (
            not self.scrolling.is_scrolling()
            and self.buttons.was_start_pressed()
        ):
            self.scrolling.set_initial_coordinates(x, y)
            self.icon.show(x, y)
            self.scrolling.start()
        elif (
            self.buttons.was_end_pressed()
            or self.buttons.was_start_released_with_hold()
        ):
            self.scrolling.stop()
            self.icon.hide()

        # it should be placed at the end to avoid initial scroll jumps
        self.scrolling.set_direction_and_coordinates(x, y)
        # debug
        self._print("click", self.debug.click)
        # start and end event have ended
        self.scrolling.clear_started_and_ended()
        # clear press information
        self.buttons.press_clear()

    def update(
        self,
        scrolling: typing.Optional[dict[str, typing.Any]] = None,
        icon: typing.Optional[ dict[str, typing.Any]] = None,
        buttons: typing.Optional[ dict[str, typing.Any]] = None,
        debug: typing.Optional[ dict[str, typing.Any]] = None,
        config: typing.Optional[ dict[str, typing.Any]] = None,
    ) -> None:
        self.config.update(**self._convert(config, {}, dict))
        self.scrolling.update(**self._convert(scrolling, {}, dict))
        self.icon.update(**self._convert(icon, {}, dict))
        self.buttons.update(**self._convert(buttons, {}, dict))
        self.debug.update(**self._convert(debug, {}, dict))

    def _is_not_end(self) -> bool:
        return not self.event_end.is_set()

    def json(self) -> dict[str, typing.Any]:
        return {
            "scrolling": self.scrolling,
            "buttons": self.buttons,
            "icon": self.icon,
            "debug": self.debug,
            "config": self.config,
        }
