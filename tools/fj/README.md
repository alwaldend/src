---
title: fj
description: Forgejo CLI
tags:
  - forgejo
---

`fj` is the repository-pinned Forgejo CLI. The
[delivery procedure](../../projects/al/skills/al-delivery/references/forgejo.md)
owns PR operations, remote selection, authentication options and capability
limits. PRs stay on the selected repository remote's forge.

On dev1, the [managed installation](../../infra/dev/ansible/README.md) provides
`~/.local/bin/fj`. With an OAuth client configured for the instance, start its
native browser login when authorized:

```sh
fj --host https://forgejo.alwaldend.com auth login
```

For the configured instance, `bazel run //tools/fj -- ...` is the existing
Vault-authenticated wrapper. To inspect the pinned CLI's help without invoking
that authentication flow, run:

```sh
bazel run //third_party/org_codeberg_forgejo_contrib_forgejo_cli_bin -- --help
```

## Links

- Site: https://forgejo.org/
- Repo: https://codeberg.org/forgejo-contrib/forgejo-cli
