---
title: Mullvad
description: Mullvad data
---

## Relays

Update relay data:

```sh
bazel run //data/mullvad:update_relays
```

## Show active servers

```sh
cat data/mullvad/relays.json  | jq -r ".wireguard.relays[] | select(.active) | .hostname"
```

{{< readfile file="relays.json" code="true" >}}
