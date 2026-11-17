# Verification

Observed 2026-11-17 against the source-only candidate based on
`cab8f25a5181cdf9de4662591e497e219a973aeb`.

- Bootstrap acceptance passes against a private TLS Vault fixture and the real
  Terraform provisioner: initial creation, immutable cache reuse, explicit
  replacement, history/access/redirect/CAS refusal, corrupt-record rejection,
  and private-value exclusion from output and state.
- The full eight-case Forgejo upgrade fixture passes with actual pinned service
  binaries. Its signing case installs the pair, makes a file-edit API commit,
  verifies it using Git and the declared public key, repairs public content and
  private permissions on a repeat without restarting, and refuses a different
  installed private key without stopping the service.
- Terraform formatting, both owning OpenSpec validations, repository quality,
  affected deployment packaging, and semantic lint pass. The final delivery
  candidate is rechecked through the same declared targets.
- Repeatable sanitized receipts are emitted by the two acceptance targets into
  Bazel's undeclared test outputs. No fixture credentials or raw state are kept
  as repository evidence.

No live plan, key generation, Vault write, or Ansible deployment was run.
Production authentication, systemd, cross-account transfers, and live service
behavior remain outside this source-only verification.

The reusable command is now tools/vault/cmd/ssh_key, with neutral mount/path
inputs. The actual provisioner acceptance fixture remains under
infra/forgejo/test/signing_key. Its relocated acceptance target passes, the new
executable path is present in setup wrapper runfiles, and all 33 selected checks
pass. The shared tool's dependency graph contains no infrastructure targets.

## Session observations

The initial 240-second suite deadline was too short; the complete run uses a
900-second bound. The SSH parser returns a pointer to an Ed25519 private key;
cache validation was corrected and rerun. The Terraform fixture must preserve
runfiles and pass the wrapper's explicit `--chdir`, rather than relying on its
process directory. Large formatter output should be redirected to task logs
from the first invocation; subsequent checks use bounded diagnostic summaries.
