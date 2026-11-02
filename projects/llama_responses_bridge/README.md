---
title: Llama Responses bridge
description: Loopback Responses API function namespace compatibility for llama.cpp
---

This project adapts Responses API function namespaces for llama.cpp servers
that accept only flat function tools. It flattens function schemas, replayed
calls, and forced tool choices, then restores namespace/name identities in JSON
and streaming replies. The adapter performs no inference or tool execution.

## Run

```sh
bazel_agent bazel run //projects/llama_responses_bridge -- --listen 127.0.0.1:8080 --upstream http://127.0.0.1:8081
```

Point the client's Responses provider at `http://127.0.0.1:8080/v1` and run
llama.cpp on `127.0.0.1:8081`. Both endpoints are configurable, but the command
requires literal IPv4 loopback addresses. It proxies only `POST /v1/responses`,
`GET /v1/models`, and `GET /health`, without query strings. The bridge has no
authentication and is intended for a trusted local host.

Request bodies, JSON replies, and individual SSE events are bounded at 32 MiB.
Only function tools and function-only namespaces are supported; other tool types
are rejected explicitly. Mappings are request-local, and disconnects cancel the
upstream stream. Call arguments and user text are not rewritten.

## Ownership and acceptance

- `cmd/responses_bridge` owns flags, loopback validation, and the HTTP process.
- `internal/responses` owns protocol translation and fixed-upstream proxying.
- `test/protocol` owns command-level fixtures with a deterministic local backend
  and an optional real Codex MCP call/result round trip.
- `openspec` owns maintained extraction work. Service units and host deployment
  remain with their infrastructure consumer.

Run the existing acceptance suite:

```sh
bazel_agent bazel test //projects/llama_responses_bridge/test/protocol:protocol_test --test_output=all
```

See [protocol acceptance](test/protocol/README.md) for the optional Codex check.
The adapter was verified with Codex 0.166.0 and llama.cpp's b11960/b12032 tool
schema limitation. These fixtures establish protocol behavior, not GPU loading,
model quality, or deployment readiness. Remove the compatibility layer when a
pinned upstream server passes the namespace round trip directly.
