---
title: Useless QT GUI
description: Desktop GUI written with C++
statuses:
  - finished
languages:
  - cpp
tags:
  - gui
  - desktop
  - qt
---

## Links

- Source code: https://github.com/alwaldend/src/tree/master/projects/useless_qt_gui

## Features

- Desktop GUI
- C++, [Qt](https://www.qt.io/)

## Build

The Bazel build downloads a SHA-256-pinned Qt 6.8.3 distribution through
`rules_qt`; a system Qt installation is not required. The project previously
targeted Qt 6.9.0. Using 6.8.3 is a deliberate downgrade to the LTS version
supported and tested by the pinned `rules_qt` release.

## TODO

- Add some screenshots
