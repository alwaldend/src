---
title: TruffleHog rules
linkTitle: TruffleHog rules
description: Scan Git history in Bazel and fail on unapproved secret findings
layout: landing
statuses:
  - in_progress
languages:
  - go
  - bzl
tags:
  - bazel
---

![Scan state machine with policy validation, finding classification, and pass, unmatched, or error outcomes](/docs/projects/rules_trufflehog/assets/scan.svg)

The [diagram key](/docs/projects/rules_trufflehog/#diagram-key) explains scan scope,
allowlist matching, and completion checks.

TruffleHog rules integrate secret scanning into Bazel tests. The rules provide a
pinned scanner and scan Git history without attempting to authenticate detected
credentials. HEAD history is the default; an explicit option includes all local refs.

Reviewed historical false positives use an exact detector, path, commit, and line
allowlist. Original findings remain available for inspection, including suppressed
ones. Inline ignores apply to the source revision being scanned, so adding an
ignore today does not suppress a finding in an older commit.

A scan passes only after successful completion with zero unmatched findings.
Malformed policy, invalid output, scanner failures, and timeouts remain errors
even when other findings have been approved.

[Read the documentation](/docs/projects/rules_trufflehog/)
