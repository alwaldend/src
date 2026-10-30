---
title: Lynis
description: Pinned CISOfy Lynis source archive for on-demand security audits
---

`//third_party/com_github_cisofy_lynis:archive` exposes the unmodified official
CISOfy source archive. It contains shell scripts, databases, profiles, and the
upstream GPLv3 license beneath `lynis/`; its executable is `lynis/lynis`.
Consumers package the archive as data and extract it before execution.

The version, download URL, and SRI pin are owned by `include.MODULE.bazel`.
The selected stable release and archive were verified on 2026-10-29. The detached
signature verified against CISOfy's published software key, primary fingerprint
`84FAA9983B24AEF24D6C87F1FEBB7D1812576482`.

- [Project and stable downloads](https://cisofy.com/downloads/lynis/)
- [Upstream releases](https://github.com/CISOfy/lynis/releases)
- [Publisher software key](https://cisofy.com/files/cisofy-software.pub)
