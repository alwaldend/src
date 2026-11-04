## 1. Implementation

- [x] 1.1 Write offline command failure fixtures before implementation; retain
      a repeatable JSON result artifact for authentication and storage boundaries.
- [x] 1.2 Replace the shell login with Go and preserve the AL entry point;
      verify the command fixtures and build the wrapped target.
- [x] 1.3 Document usage, native transport requirements, and validation limits;
      verify help output against the README.

## 2. Validation

- [x] 2.1 Run strict OpenSpec validation, repository formatting/quality, affected
      semantic lint, and the command tests against the final candidate.
      Git and the AL preparation/publication receipt own delivery state. Publication
      uses the requested `gh` transport, with the receipt retained under
      `out/sort-projects-by-creation-date/`.

## Evidence

The command fixture passes 20 scenarios and writes `login-results.json` under
its Bazel undeclared test outputs. The AL-wrapped `//tools/vault/login` target
builds. A read-only smoke check using the real hardware's public certificate
successfully completed automatic selection before transport was directed to
an unused loopback port. It neither authenticated to Vault nor created a token
file. The operator subsequently confirmed a successful live Vault login after
the private-key selector was corrected to match by object ID instead of the
certificate's label. Authenticated read-only access to GitLab then succeeded.
