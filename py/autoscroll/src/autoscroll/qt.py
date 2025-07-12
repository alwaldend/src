from sys import argv as sys_argv

from PySide6.QtCore import Qt
from PySide6.QtSvgWidgets import QSvgWidget
from PySide6.QtWidgets import QApplication

from .functions import get_path

application = QApplication(sys_argv)
application.setQuitOnLastWindowClosed(False)


class Icon(QSvgWidget):

    def __init__(self, path: str, size: int) -> None:
        super().__init__()
        self.update_icon(path, size)
        self.renderer().setAspectRatioMode(Qt.KeepAspectRatio)
        self.setWindowFlags(
            Qt.FramelessWindowHint
            | Qt.WindowStaysOnTopHint
            | Qt.WindowTransparentForInput
            | Qt.Tool
        )
        self.setStyleSheet("background:transparent")
        self.setAttribute(Qt.WA_TranslucentBackground)

    def update_icon(self, path: str, size: int) -> None:
        self.load(get_path(path))
        self.icon_size = size
        self.resize(size, size)

    def show(self, x_current: int, y_current: int) -> None:
        half_size = self.icon_size // 2
        self.move(x_current - half_size, y_current - half_size)
        super().show()
