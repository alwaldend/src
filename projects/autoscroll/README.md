---
title: Autoscroll
description: CLI app for autoscroll
tags:
  - py
---

## About The Project

This package implements platfrom-independent autoscroll

## Features

- Pretty useless
- Has a config file with hot-reload
- Python3, PyQt6
- Some argparse tinkering

## Usage

By default, the icon is disabled, to enable it pass `--icon-enable`

You can pass file contents as command line arguments using `@path/to/the/file` syntax.
Arguments in that case can be placed wherever - on one line, on several lines

If you want to dynamically pass runtime arguments (without restarting the process), you can use `--config` options for it

Once you press `--buttons-start`, you can scroll vertically or horizontally just by moving your mouse untill you press `--buttons-end`

If `--buttons-hold` is set, the srolling ends once you release `--buttons-start`

Once `--buttons-start` is pressed, the scroll thread starts looping
Every loop consists of sleeping for an interval, then scrolling for either 0, 1, or -1 pixels on both axis towards the starting point
Starting point is the point where `--buttons-start` was pressed
Sleep interval is recalculated on every mouse move as such:

```
    100 / (--scrolling-acceleration * max(distance) + --scrolling-speed)
```

If `--scrolling-acceleration` is not 0, the speed of scrolling will be faster
the farther away you are from the starting point
If `--scrolling-acceleration` is 0, the speed of scrolling will be constant

### Examples

#### Use the package

```bash
python3 -m venv venv
. venv/bin/activate
pip install autoscroll
autoscroll
```

#### Command line options

```bash
autoscroll --buttons-start 1 --debug-click --icon-disable
```

#### Pass a configuration file once

```bash
autoscroll --icon-enable @config.txt
```

If `config.txt` is defined like this, its contents will be used as command line arguments - they will be loaded only once
Arguments can be placed wherever - on one line, on several lines
For example,

```
--buttons-start 1
--buttons-hold --debug-click
```

#### Listen for changes in the configuration file

```bash
autoscroll --config-enable --config-path config.txt
```

If config.txt is defined like this, the process will listen for changes in that
file and update itself
Arguments can be placed wherever - on one line, on several lines
The file is checked for changess every `--config-interval`
For example:

```
--buttons-start 1 --buttons-hold
--debug_click
```

### `--help` output

```
usage: autoscroll [-h] [-ss SCROLLING_SPEED] [-sd SCROLLING_DEAD_AREA]
                  [-sa SCROLLING_ACCELERATION] [-bh] [-bs BUTTONS_START] [-be BUTTONS_END]
                  [-ce] [-cp CONFIG_PATH] [-ci CONFIG_INTERVAL] [-ie] [-ip ICON_PATH]
                  [-is ICON_SIZE] [-df] [-dc] [-ds] [-di]

...

options:
  -h, --help            show this help message and exit

scrolling:

  -ss, --scrolling-speed int
                        constant part of the scrolling speed
                        [default: 300]
  -sd, --scrolling-dead-area int
                        size of the square area aroung the starting point where scrolling will stop, in
                        pixels
                        [default: 50]
  -sa, --scrolling-acceleration int
                        dynamic part of the scrolling speed, depends on the distance from the point
                        where the scrolling started, can be set to 0
                        [default: 10]

buttons:

  -bh, --buttons-hold   if set, the scrolling will end once you release --buttons-start
  -bs, --buttons-start int
                        button that starts the scrolling
                        [default: 2]
  -be, --buttons-end int
                        button that ends the scrolling
                        [default: --buttons-start]

config:

  -ce, --config-enable  if set, arguments from the configuration file on --config-path will be loaded
                        every --config-interval
  -cp, --config-path str
                        path to the configuration file
                        [default: ~/.config/autoscroll/config.txt]
  -ci, --config-interval int
                        how often the config file should be checked for changes, in seconds
                        [default: 5]

icon:

  -ie, --icon-enable    if set, the icon will be enabled
  -ip, --icon-path str  path to the icon
                        [default: resources/img/icon.svg]
  -is, --icon-size int  size of the icon, in pixels
                        [default: 30]

debug:

  -df, --debug-file     if set, every time the config file is parsed, information will be printed to
                        stdout
  -dc, --debug-click    if set, click info will be printed to stdout
  -ds, --debug-scroll   if set, scroll info will be printed to stdout
  -di, --debug-initial  if set, startup configuration will be printed to stdout
```
