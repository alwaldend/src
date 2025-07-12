from threading import Event, Thread
from time import sleep
from typing import Any, Dict

from pynput.mouse import Button, Listener

from .support import Base, Buttons, Config, Debug, Icon, Scrolling


class Autoscroll(Base):

    def __init__(self, *args, **kwargs) -> None:
        self.scrolling: Scrolling = Scrolling()
        self.icon: Icon = Icon()
        self.config: Config = Config()
        self.buttons: Buttons = Buttons()
        self.debug: Debug = Debug()
        self.event_end: Event = Event()

        # update from initializer arguments
        self.update(*args, **kwargs)

        # threads
        # mouse actions
        self.thread_scroll_listener = Listener(
            on_move=self._on_move, on_click=self._on_click, daemon=True
        )
        # scroll
        self.thread_scroll_action = Thread(
            target=self._loop,
            daemon=True,
            args=(self._is_not_end, self._scroll),
        )
        # update from the config file (if enabled)
        self.thread_config = Thread(
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
        Thread(
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
        # wait untill the config file is enabled
        self.config.wait()
        # update from the config file
        self.update(**self.config.parse_config_file())
        # debug
        self.config._print("config", self.debug.file, ["content"])
        # sleep
        sleep(self.config.interval)

    def _scroll(self) -> None:
        # wait for the scrolling event to be set in _on_click
        self.scrolling.wait()
        # wait for a period of time, calculated in _on_move
        self.scrolling.sleep_for_interval()
        # scroll on x-axis and y-axis, each scroll is either 1px, 0px, or -1px
        self.scrolling.scroll_once()

    def _on_move(self, x: int, y: int) -> None:
        # update coordinates
        self.scrolling.set_direction_and_coordinates(x, y)
        # recalculate sleep interval
        self.scrolling.set_interval()

    def _on_click(self, x: int, y: int, button: Button, pressed: bool) -> None:
        # send information about which button was pressed/released
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
        scrolling: Dict[str, Any] = None,
        icon: Dict[str, Any] = None,
        buttons: Dict[str, Any] = None,
        debug: Dict[str, Any] = None,
        config: Dict[str, Any] = None,
    ) -> None:
        self.config.update(**self._convert(config, {}, dict))
        self.scrolling.update(**self._convert(scrolling, {}, dict))
        self.icon.update(**self._convert(icon, {}, dict))
        self.buttons.update(**self._convert(buttons, {}, dict))
        self.debug.update(**self._convert(debug, {}, dict))

    def _is_not_end(self) -> bool:
        return not self.event_end.is_set()

    def json(self) -> Dict[str, Any]:
        return {
            "scrolling": self.scrolling,
            "buttons": self.buttons,
            "icon": self.icon,
            "debug": self.debug,
            "config": self.config,
        }
