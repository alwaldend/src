---
title: OpenHands canvas UI compatibility module
description: Pinned OpenHands Agent Canvas compatibility tool
tags:
  - third_party
  - openhands
---

`tools/canvas_ui_tool.py` from `OpenHands/OpenHands` at tag `v1.18.0`
(commit `9120ff6cbbe23640f0e475661e5a9c9729cdbf1f`). The OpenHands Agent
Canvas launchers import it into the agent server through
`--import-modules canvas_ui_tool`.

It is required because the automation preset builds its agent with a
`finish_tool_response_schema`, which registers the SDK's builtin `FinishTool`
only inside the preset's own process while advertising it to the agent server as
`openhands.sdk.tool.builtins.finish`. That SDK module does not self-register, so
every automation-dispatched remote conversation fails to resolve `FinishTool`
unless the compatibility module registers it in the agent server process first.

The module also keeps persisted pre-`client_tools` conversations resolvable.
It is vendored rather than copied so the upstream revision, license, and
provenance stay with the third-party tree.

Upstream license: MIT, see
[LICENSE](https://github.com/OpenHands/OpenHands/blob/v1.18.0/LICENSE).
