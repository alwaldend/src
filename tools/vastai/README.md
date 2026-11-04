---
title: Vast.ai CLI
description: Pinned upstream command-line client for Vast.ai GPU cloud resources
---

Runs the unmodified upstream `vastai` CLI, pinned to 1.10.0 in
[`pyproject.toml`](pyproject.toml). The existing
[`Python dependency workflow`](../py/README.md) defines the update mechanism: `bazel run //tools/vastai:uv_lock.update`.
The component-local lock isolates upstream's exact dependency pins from the shared
Python environment.

```sh
bazel run //tools/vastai -- --help
bazel run //tools/vastai -- search offers --help
```

This target does not install software on the host or supply credentials.
Commands that contact Vast.ai require network access; resource operations can
incur charges. Configure authentication through upstream's supported credential
mechanism before making authenticated requests. Keep API keys out of source,
command arguments, and captured output.

## Links

- [Upstream source and usage](https://github.com/vast-ai/vast-cli)
- [Pinned release](https://github.com/vast-ai/vast-cli/releases/tag/v1.10.0)
