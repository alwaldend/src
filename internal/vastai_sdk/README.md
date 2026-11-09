---
title: Vast.ai SDK
description: A bounded client backed by the official generated API protocol
---

The SDK exposes one public interface at `pkg/client` for offers, instances,
SSH public keys, environment variables, team roles, teams, project API keys,
and exact template lookup. Template reads validate the returned hash, image,
and Docker option metadata; they return an environment-only classification
without retaining environment values. Creation accepts an optional template hash and Docker environment/port options.
Instance reads expose typed host-port bindings without retaining returned environment values.
The Terraform provider consumes this client; it owns Terraform schemas and state.

Construct a client with `client.New(client.Options{Endpoint: endpoint, APIKey: key})`.
Both fields are required. Role names remain individual URL segments, including
names containing an encoded slash.

`internal/generated` supplies protocol operations and models generated offline
from the [official schema pin](../../third_party/ai_vast_openapi/README.md).
Generated source stays in Bazel outputs. `generator.json` owns configuration;
the existing OpenAPI Generator is pinned to 7.26.0.

The public client adds 60-second deadlines, HTTPS endpoint validation, an 8 MiB
response bound, redirect refusal, and credential-safe errors. It checks endpoint
success and resource identity, including the distinction between omitted and
explicit-null instance responses. Those lifecycle guarantees are absent from
the raw generated interface, which remains an implementation detail.

The schema omits account SSH-key listing and the console's context-selection
endpoint; the client supplies those operations. API-key reads use the generated
key-collection operation. Revocation checks that collection before and after
deletion, so replies without success flags do not abandon cleanup ownership.
Team-context selection verifies the current account first and leaves an
already-selected team key unchanged. Canonical trailing-slash API paths avoid
redirects; redirect refusal remains enabled.
The search request maps to an open JSON object to preserve arbitrary filters
and avoid invalid object-default initializers in the pinned generator. Numeric
mappings retain 64-bit identifiers and Terraform-compatible precision.

```sh
bazel build //internal/vastai_sdk/pkg/client
```

The [provider integration](https://github.com/alwaldend/src/pull/249) exercises
this client through real Terraform CLI lifecycle fixtures against loopback API
servers. Live mutation compatibility remains unverified. This internal library
has no independent product release or landing page.
