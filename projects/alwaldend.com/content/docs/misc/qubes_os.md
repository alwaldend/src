---
title: Qubes OS
decription: Qubes OS
tags:
  - qubes_os
---

## Links

- Site: https://www.qubes-os.org/

## Restore normal desktop zoom

If the whole desktop is magnified and follows the mouse pointer in Xfce,
move the pointer over the desktop or taskbar, hold **Alt**, and **scroll
down** until the screen returns to normal size. Holding **Alt** while
scrolling up zooms in.

See the [Xfce desktop zoom discussion](https://forum.xfce.org/viewtopic.php?id=17905).

## Qube startup timeout

Run in dom0:

```sh
qmv-prefs <VM NAME> qrexec_timeout <seconds>
```
