# Useless QT GUI

## Purpose

Provide a C++ desktop application with Qt widgets and custom menu and window
controls. This baseline was observed at repository revision `550d7e79` on
2026-09-08. Sources are the
[project README](../../../README.md),
[build definition](../../../BUILD.bazel),
[entry point](../../../main/cpp/main.cpp),
[main window](../../../main/cpp/mainwindow.cpp), and
[menu controls](../../../main/cpp/menubutton.cpp).

## Requirements

### Requirement: Build the application through Qt Bazel rules

The project SHALL build its application through `rules_qt`, compile the
`mainwindow.ui` form, package the QRC resources, and link the application with
Qt Widgets. The build SHALL use the repository's pinned Qt distribution
without requiring a system Qt installation.

#### Scenario: Build application target

- **WHEN** Bazel builds `//projects/useless_qt_gui:useless_qt_gui`
- **THEN** the target includes the application library, generated main-window UI, packaged resources, and Qt Widgets dependency

### Requirement: Open a custom desktop window

The executable SHALL construct a `QApplication`, show its main window, and run
the Qt event loop. The main window SHALL provide a menu area, a stacked content
area, and controls for minimizing, maximizing, and closing the application.

#### Scenario: Launch the executable

- **WHEN** the executable starts in a working graphical environment
- **THEN** it displays the frameless main window and enters the Qt event loop

### Requirement: Add and select content pages

The menu SHALL provide unused add controls that become named content-page
buttons when activated. Selecting a content-page button SHALL select its
associated page in the stacked content area.

#### Scenario: Activate an unused content slot

- **WHEN** an unused add control of type `6` is activated
- **THEN** the window creates a content page, assigns its index to that control, and converts the control to a content-page button

#### Scenario: Select an existing content page

- **WHEN** an existing content-page button is activated
- **THEN** the window sets the content stack's current index to the button's page index
