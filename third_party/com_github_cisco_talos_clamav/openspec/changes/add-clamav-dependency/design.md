## Context

See proposal.md for the approved dependency scope. The official DEB is an ar
archive containing `data.tar.gz`. Bazel's existing `http_archive` supports DEB
extraction, and consumers can unpack the inner tar directly.

## Goals / Non-Goals

Expose upstream bytes without installing the package, rewriting ELF files, or
introducing another extraction dependency. Live scanning, configuration, and
signature updates remain with the consumer.

## Decisions

- Use `http_file` for the original DEB and `http_archive` for its payload tar.
  Both consume one URL and integrity definition in the owning module include.
- Export the unchanged `data.tar.gz` through `payload_tar`. The existing
  `al_pkg_extract_dir` accepts only `.tar`; a new directory rule or repackaging
  adds no value to the consumer's accepted tar interface.
- Retain upstream paths, symlinks, license files, and bundled libraries.
  Consumers set `LD_LIBRARY_PATH` to their extracted `usr/local/lib` because
  the executables contain an absolute upstream library search path.

## Risks / Trade-offs

- Wrong artifact or corruption: enforce SRI and compare the downloaded hash
  and size to publisher release metadata.
- Unsupported CPU or libc: document Linux x86_64 with glibc 2.30 or newer;
  this package is for the controller, not XCP endpoints.
- Extraction loses modes or symlinks: compare the Bazel tar byte-for-byte
  with the DEB member and run extracted `clamscan`, `freshclam`, and `clamd`
  with `--version` in ignored task scratch.
- Package resolution must not run installer scripts or services: use native
  archive extraction only. No database update or scan is part of validation.
- Detached signature verification is not yet established; report that limit
  separately from the checked SHA-256 and publisher metadata match.
