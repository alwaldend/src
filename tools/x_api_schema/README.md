---
title: X API schema fetch tool
description: Authenticated upstream schema download for dependency maintenance
---

This first-party developer tool implements the X API schema download command.
Its maintained entry point belongs to the dependency package:

```sh
bazel_agent bazel run //third_party/com_x_api_openapi:fetch_schema
```

That target injects the Article uploader's existing OAuth 1.0a credentials.
The command downloads into ignored `out/x_sdk/openapi.json` by default and
accepts `--out` to select another destination. It validates a bounded response
and atomically replaces the download only after success. It does not publish
the snapshot or alter dependency pins. The dependency's
[update procedure](../../third_party/com_x_api_openapi/README.md) owns those steps.

The command and its local HTTP E2E fixtures live in `cmd/fetch_schema`.
The [runtime SDK](../../internal/x_sdk/README.md) consumes the pinned mirror and
does not own this maintenance command.
