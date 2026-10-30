---
title: ClamAV
description: Pinned upstream ClamAV payload for Linux x86_64 controllers
---

`//third_party/com_github_cisco_talos_clamav:archive` exposes the unmodified
official DEB. `:payload_tar` exposes its unchanged `data.tar.gz`, extracted by
Bazel's native DEB support. The tar retains `usr/local/bin`, `usr/local/sbin`,
`usr/local/lib`, sample configuration, headers, documentation, and licenses.
Consumers can unpack it with Ansible `unarchive` or tar without stripping paths.
Dependency resolution runs no package installation scripts or services.

The version, URL, and SRI integrity are owned by `include.MODULE.bazel`.
The selected stable release, downloaded SHA-256, and byte size were checked
against the publisher's GitHub release metadata on 2026-10-29. The detached
signature has not been verified.

## Controller runtime

The payload requires Linux x86_64 with glibc 2.30 or newer. It is intended for
the scan controller; it is not compatible with older XCP endpoint runtimes.
Consumers retain the bundled libraries and explicitly set `LD_LIBRARY_PATH`
to `<extraction-root>/usr/local/lib` because the upstream executables contain
an absolute library search path. For example, after unpacking into task scratch:

```sh
clamav_root="$PWD/out/host-security-updates-20261029/clamav/dependency/payload"
LD_LIBRARY_PATH="$clamav_root/usr/local/lib" \
  "$clamav_root/usr/local/bin/clamscan" --version
```

Consumers own explicit configuration paths, signature databases, database
updates, and scan commands. The dependency does not supply a populated database.
`freshclam --version` also reads configuration; pass `--config-file` explicitly
even for a version check. An empty file suffices for that check only.
Keep upstream copyright and license files when packaging or deploying it;
ClamAV is distributed under GPL-2.0, with separately licensed bundled components.

## Upstream

- [ClamAV releases](https://github.com/Cisco-Talos/clamav/releases)
- [Selected release](https://github.com/Cisco-Talos/clamav/releases/tag/clamav-1.5.4)
- [Installation documentation](https://docs.clamav.net/manual/Installing.html)
- [Upstream license](https://github.com/Cisco-Talos/clamav/blob/clamav-1.5.4/COPYING)
