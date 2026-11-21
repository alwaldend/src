## Why

Normal CI covers formatting and several generator fixtures. Generated BUILD
file drift and the semantic linter profile need explicit CI coverage.

## What Changes

- Check root-workspace Gazelle output without rewriting the checkout.
- Run the repository semantic lint profile in CI and report its result.
- Correct the direct archive declarations that warn during module resolution.

## Capabilities

- `repository` gains CI gates for BUILD generation and semantic lint.
