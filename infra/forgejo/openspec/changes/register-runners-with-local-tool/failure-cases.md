## Acceptance cases defined before implementation

- Missing credentials or invalid endpoint/path: reject before network mutation.
- Vault access denied or invalid stored record: fail without runner creation.
- Existing matching cached registration: reuse without a Forgejo request or Vault write.
- Cached record remains even when a runner is externally removed: preserve
  it without automatic invalidation or rotation.
- Existing same-name runner on a later page: reject without adopting or deleting.
- Concurrent applies use the owning Terraform backend lock; independent
  concurrent command invocation is outside the command contract.
- Creation error or lost response: a same-name runner blocks duplicate creation.
- Final Vault write failure after successful creation: same-name collision
  stops the retry without duplicate creation.
- Redirect or untrusted TLS endpoint: reject without forwarding credentials.
- Upstream errors echo credentials: command output contains no credential values.
- Terraform plan and repeat apply: no creation; initial isolated apply: one POST;
  explicit replacement: verify ready registration; removal: no automatic deletion.

Fixtures use only public dummy credentials and private loopback test services.
The repeatable artifact records case outcomes without secret values.
