## Context

The AL wrapper injects `VAULT_ADDR` and `VAULT_CACERT` without authenticating.
The old shell script discovers the last token, calls curl, parses JSON with jq,
and truncates `~/.vault-token` directly. The host's curl uses OpenSSL 4; its
installed PKCS#11 provider loads successfully, but engine support is absent.

## Goals / Non-Goals

Keep the public Bazel login target and role-name argument. Replace the shell
pipeline with a Cobra Go command and prevent token loss on failed login.
Live authentication and persistent host configuration are separate operations.

## Decisions

Proceed with Go orchestration and the existing curl PKCS#11 TLS boundary. Go's
standard TLS package does not load PKCS#11 certificates or signers; using native
Go TLS would require a new hardware integration or cryptography dependency.
A new dependency is unnecessary for replacing the script. Curl receives explicit
certificate/key types and provider selection, without a shell or user curlrc.
Go supplies JSON through stdin, bounds output, validates the response, and uses
a same-directory private temporary file and atomic rename for persistence.

Replace arbitrary trailing curl arguments with typed certificate, token-file,
mode, and timeout options. Reject plaintext or credential-bearing Vault URLs,
PIN-bearing certificate URIs and symlink token files. Automatic discovery ignores
trust modules, reads public hardware certificates, filters validity/client usage,
and prefers an identity matching the role name. The newest matching certificate
is selected; an unmatched sole usable certificate supports role aliases. An
explicit URI remains an optional override. Multiple unmatched client identities
require explicit selection because discovery cannot infer Vault's trust policy.
Private-key selection preserves the certificate ID and removes its object label:
PKCS#11 certificates and private keys may have different labels for the same ID.
Response bodies and certificate-discovery diagnostics remain private. Curl's
normal TLS diagnostics and provider prompts pass through stderr so certificate
loading failures can be diagnosed. Verbose HTTP logging is not enabled.

## Risks / Trade-offs

Offline command fixtures verify provider selection and failure boundaries but
cannot establish a real YubiKey login. Native curl and p11tool remain host
prerequisites. Provider mode is the default for the observed host; explicit
engine mode supports older installations without changing global OpenSSL config.

## Failure cases before implementation

Missing role/address/CA, insecure or malformed URLs, PIN-bearing URIs, absent
tokens, discovery failure, curl failure or timeout, malformed or
oversized JSON, missing/invalid token, unwritable destination, symlink target,
and failed replacement must preserve the previous token and hide secret data.

On this host, read-only discovery observed two trust modules and one hardware
token. The hardware certificate is currently valid, has TLS client-authentication
usage, and matches the checked-in Vault certificate role. No real login was
performed to collect this evidence.
