---
title: Hugo
links:
  - url: https://gohugo.io
---

Hugo is a static site builder

## Environment variables

Hugo has an allowlist of environment variables, and `js_binary` rules
do not work because they need a `BAZEL_BINDIR` variable

Fix:

```toml
[security.exec]
osEnv = [
    '(?i)^((HTTPS?|NO)_PROXY|PATH(EXT)?|APPDATA|TE?MP|TERM|GO\w+|(XDG_CONFIG_)?HOME|USERPROFILE|SSH_AUTH_SOCK|DISPLAY|LANG|SYSTEMDRIVE|BAZEL.+)$',
]
```

Configuration reference: https://gohugo.io/configuration/security/
