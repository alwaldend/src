---
title: Ceph
description: Ceph
tags:
  - ceph
---

## Links

- Proxmox docs: https://pve.proxmox.com/wiki/Deploy_Hyper-Converged_Ceph_Cluster
- Single node crushmap: https://imanudin.net/2025/05/12/proxmox-ceph-single-node-installation-for-test-lab-or-home-setup/

## See current crush map

Go to Nodes -> Host -> Ceph -> Configuration

## Crush map update

Ssh to the PVE host and update the map:
```sh
sudo ceph osd getcrushmap -o crushmap.cm
sudo crushtool --decompile crushmap.cm -o crushmap.txt
sudo vim crushmap.txt
sudo crushtool --compile crushmap.txt -o new_crushmap.cm
sudo ceph osd setcrushmap -i new_crushmap.cm
sudo ceph -s
```

