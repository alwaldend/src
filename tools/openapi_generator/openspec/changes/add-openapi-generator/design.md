## Context

See proposal.md. The repository already has Java and Go toolchains.

## Goals / Non-Goals

Keep external pins in third_party and generated source consumable without execution tools. No infrastructure or secret changes are required.

## Decisions

Pin the self-contained 7.26.0 Maven Central JAR using http_jar and SRI. This uses the existing Java toolchain and avoids a new package manager or transitive Maven lock. The generation rule produces a Bazel tree artifact. Consumers filter the generated source they need and compile it without writing build artifacts to source.

## Risks / Trade-offs

Upstream schema and generated names may change; compile and behavior checks validate each refreshed candidate. Authentication availability may prevent live fetching; retain that failure explicitly.
