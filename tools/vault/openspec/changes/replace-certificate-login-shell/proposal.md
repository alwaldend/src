## Why

The certificate login shell pipeline relies on implicit OpenSSL engine selection
and fails on this host's OpenSSL 4. Replace it with a Go command that validates
inputs and preserves the existing token when authentication fails.

## What Changes

- Keep `//tools/vault/login` as the AL-wrapped entry point, backed by Go.
- Use explicit PKCS#11 provider certificate and key types by default.
- Parse discovery and login responses in Go, with bounded subprocess output.
- Atomically save a successful token with owner-only permissions.
- Add offline command-level verification with synthetic hardware/TLS tools.

## Capabilities

### New Capabilities

- `certificate-login`: Authenticate with a PKCS#11 client certificate and
  safely replace the host token file.

### Modified Capabilities

None.

## Impact

`tools/vault/login` source, build declarations, documentation, and command tests.
Curl remains the existing PKCS#11 TLS transport; Go owns orchestration and
storage. No new external dependency or persistent host configuration is needed.
The Hugo ordering task remains active and awaits historical date verification.
