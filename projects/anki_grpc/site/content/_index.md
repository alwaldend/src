---
title: Anki gRPC
linkTitle: Anki gRPC
description: Use Anki's native backend through gRPC
layout: landing
statuses:
  - finished
languages:
  - rust
tags:
  - anki
  - grpc
---

Anki gRPC exposes Anki's native backend to applications in other languages.
Clients use Anki's own protobuf messages to open collections, edit notes,
and invoke backend operations without running the desktop application.

The server is a Rust CLI with stdin/stdout as its default transport. TCP and
Unix sockets are available through command-line options. The backend is ready
at startup, and Anki owns collection behavior. Applications own server lifetime
and recovery after failures.

Native builds currently support Linux x86_64.

[Read the documentation](/docs/projects/anki_grpc/)
