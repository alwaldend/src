---
title: Ping
---

## `sendmsg: operation not permitted`

This problem can be caused by vpns with "kill-switch" enabled, so you might need to disable them. 

Disable Mullvad VPN:

```sh
sudo systemctl disable --now mullvad-daemon
```

