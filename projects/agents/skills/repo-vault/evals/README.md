---
title: Repository Vault evaluations
---

# Repository Vault evaluations

This suite describes the safe token-refresh contract for Vault operations on
this host. Its required offline Bazel target validates the Promptfoo
configuration, referenced case, and staged skill without making a model call.

A live target is omitted because representative behavior requires a reachable
Vault server and valid host credentials. Offline validation does not prove a
token can be refreshed; that requires the issuing client certificate to be
present and Vault to be unsealed.
