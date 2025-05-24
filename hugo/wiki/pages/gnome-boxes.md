---
title: Gnome boxes
---

[Gnome boxes](https://apps.gnome.org/Boxes/) is a hypervisor, it works decently.

## Cpu

The default cpu doesn't have a lot of capabilities, so you might need to patch
the cpu config co compile some stuff:

```xml
  <!-- https://www.qemu.org/docs/master/system/i386/cpu.html -->
  <cpu mode="host-model" />
```
