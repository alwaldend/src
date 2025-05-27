---
title: Shell
---

## Deactivate venv in subshell

```sh
. .venv/bin/activate
(
    echo "subshell, still in venv: $(which python)" && \
    . "$(dirname "$(which python)")/activate" && \
    deactivate && \
    echo "subshell, out of venv: $(which python)"
)
```

