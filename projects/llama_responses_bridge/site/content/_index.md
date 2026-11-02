---
title: Llama Responses bridge
linkTitle: Llama Responses bridge
description: Local function namespace compatibility between Codex and llama.cpp
layout: landing
statuses:
  - active
languages:
  - go
tags:
  - ai
  - cli
---

Llama Responses bridge connects clients that send Responses API function
namespaces, such as Codex, to llama.cpp servers that accept flat function tools.
It preserves tool identities across requests, streamed replies, and tool-result
follow-ups, so the client can execute the intended tool.

## How it works

Run the bridge between your client and a local llama.cpp server. Point the
client at the bridge's `/v1` endpoint; the bridge translates function namespaces
and forwards inference to llama.cpp. Tool execution stays with the client.

Both connections are restricted to IPv4 loopback. The bridge has no
authentication and is intended for a trusted local host. It supports JSON
and streaming Responses replies and rejects unsupported tool types explicitly.

## Status

The adapter and command-level protocol fixtures are implemented. GPU model
loading, model quality, and deployment readiness require separate validation.
The bridge is a compatibility layer for server versions without native namespace
support; it can be removed when the server passes the tool-call round trip directly.

[Documentation](/docs/projects/llama_responses_bridge/) ·
[Source](https://github.com/alwaldend/src/tree/master/projects/llama_responses_bridge)
