---
title: Gnome boxes
---

[Gnome boxes](https://apps.gnome.org/Boxes/) is a hypervisor.

## Cpu

The default cpu doesn't have a lot of capabilities, so you might need to patch
the config:

```xml
  <!-- https://www.qemu.org/docs/master/system/i386/cpu.html -->
  <cpu mode="host-model" />
```

### Network

Default libvirt network interferes with gnome-boxes, you might need to disable it:

```sh
sudo systemctl disable --now libvirtd
```

### Disk size

Just increasing the disk size is not enough, you need to increase the logical volume
size:

```sh
sudo growpart /dev/vda 3
sudo lvextend -l +100%FREE /dev/fedora/root
sudo xfs_growfs /dev/fedora/root
```

Links:

- https://www.redhat.com/en/blog/resize-lvm-simple
- https://askubuntu.com/questions/1406697/extend-lvm-size
