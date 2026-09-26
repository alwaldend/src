---
title: Yandex Cloud
description: Static-hosting VM with retained local content storage
---

Provisions one Fedora VM in its component folder: two CPUs, 2 GiB RAM, a
20 GiB boot disk, and a separate 100 GiB content disk. The disk and reserved
IPv4 address have destruction guards and survive replacement of the VM.
The instance firewall permits TCP 22/80/443; Nginx has no public listener.

The delegated Yandex DNS zone's A record follows the reserved address.
Cloudflare delegation and service aliases belong to the sibling DNS root.
AL injects the folder ID and folder-scoped service account; no credentials
or allocated addresses are copied into checked-in variables.
