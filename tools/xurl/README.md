---
title: X CLI
description: Pinned upstream xurl command-line client for the X API
tags:
  - x
  - cli
---

This target runs the unmodified upstream xurl CLI. The release archive and
integrity pin belong to
[`third_party/com_github_xdevplatform_xurl`](../../third_party/com_github_xdevplatform_xurl/README.md).
The current toolchain supports Linux x86-64.

```sh
bazel_agent bazel run //tools/xurl -- version
bazel_agent bazel run //tools/xurl -- --help
```

xurl uses its native `~/.xurl/auth.yml` credential store. Startup can create or
migrate this store and import `~/.twurlrc`, even for help and version commands.
There is no upstream configuration-path override. Automated use must isolate
these paths; do not place credentials in command arguments or enable verbose
logging, which prints the Authorization header.

Select `--auth oauth1` and `--app <name>` explicitly when comparing the Article
uploader with the CLI. Otherwise xurl prefers stored OAuth 2.0 credentials.
`--data` accepts literal JSON, not curl-style `@file` input. The CLI follows
redirects and uses Go's default transport; a single invocation does not provide
the uploader's stricter redirect and request-replay safeguards.

This generic tool supplies neither repository credentials nor a Markdown
converter. The [Article uploader](../../projects/x_article_uploader/README.md)
owns those operations and its Vault injection flow. Creating a draft is a
separate, quota-consuming operation from inspecting this CLI.

## Links

- [Upstream source](https://github.com/xdevplatform/xurl)
- [X CLI documentation](https://docs.x.com/tools/xurl)
