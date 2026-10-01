---
title: X Article upload skill evaluations
---

The offline target validates that Promptfoo can load the skill, its cases, and
the complete staged configuration without credentials or model calls. The cases
cover the convert-then-draft sequence and the unpublished default, reading the
24-hour per-user cap from the response headers, keeping the short endpoint
window separate from that cap, distinguishing an expired Vault bootstrap token
from an X rejection, and reusing the cached media identifier after a capped
draft call.

A live target is omitted because realistic evaluation requires Vault-injected
credentials and network calls that spend a limited daily draft allowance.
Configuration validation does not establish that an upload succeeds or that a
draft exists; live evidence belongs to the actual authorized upload.
