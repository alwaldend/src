## Why

The Responses namespace adapter is reusable software and should be reviewed
independently of the V100 host deployment.

## What Changes

- Extract the existing bridge, internal protocol adapter, and command-level
  acceptance suite into `projects/llama_responses_bridge`.
- Give the project its own documentation, Bazel targets, OpenSpec owner,
  and registered landing page in the shared site.
- Retain loopback restrictions, function identities, and protocol behavior.

## Capabilities

No spec-level behavior changes. This is an ownership and packaging refactor;
`skip_specs: true` records that distinction.

## Impact

The extracted project uses the existing Go and Cobra dependencies. The
separate local inference PR will consume the command target. No host deployment
or new external dependency is included.
