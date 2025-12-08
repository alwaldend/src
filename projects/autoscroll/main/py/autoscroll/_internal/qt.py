from sys import argv as sys_argv

import PySide6.QtCore
import PySide6.QtSvgWidgets
import PySide6.QtWidgets

from .functions import get_path

application = PySide6.QtWidgets.QApplication(sys_argv)
application.setQuitOnLastWindowClosed(False)

class Icon(PySide6.QtSvgWidgets.QSvgWidget):

    def __init__(self, path: str, size: int) -> None:
        super().__init__()
        self.update_icon(path, size)
        self.renderer().setAspectRatioMode(PySide6.QtCore.Qt.KeepAspectRatio)
        self.setWindowFlags(
            PySide6.QtCore.Qt.FramelessWindowHint
            | PySide6.QtCore.Qt.WindowStaysOnTopHint
            | PySide6.QtCore.Qt.WindowTransparentForInput
            | PySide6.QtCore.Qt.Tool
        )
        self.setStyleSheet("background:transparent")
        self.setAttribute(PySide6.QtCore.Qt.WA_TranslucentBackground)

    def update_icon(self, path: str, size: int) -> None:
        self.load(get_path(path))
        self.icon_size = size
        self.resize(size, size)

    def show(self, x_current: int, y_current: int) -> None:
        half_size = self.icon_size // 2
        self.move(x_current - half_size, y_current - half_size)
        super().show()
