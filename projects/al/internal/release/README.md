---
title: Release operations
description: Artifact publication, website operations, preview, and inspection
---

This package owns artifact identity, the project publication lifecycle, and the
website extraction and activation operations. See
the [AL release documentation](../../README.md#preview-and-inspect-artifact-publication)
for operator commands and their guarantees.

`al release publish` uploads artifacts; `al release website extract` creates a retained
site and `al release website activate` selects it. Each website command supports offline
preview with `--dry-run`; `al release plan` previews upload. Mutating commands
can checkpoint `--receipt`, while a website dry run rejects that flag.

These operations share destination selection and canonical manifest preparation.
Read-only `al release inspect` keeps local payload
availability independent of the remote snapshot, so operators can investigate
an uncertain publication after local build outputs have been removed.

The observation transport uses the existing strict SSH settings and content
account. Its remote script performs reads only, without locks or staging. Each
captured stream is limited to 16 KiB; overflow cancels the subprocess. Unsafe
filesystem state is reported without repair, while transport and read failures
remain unavailable evidence. A snapshot is not a lease against later changes.
