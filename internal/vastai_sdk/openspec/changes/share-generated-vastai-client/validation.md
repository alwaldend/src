# Validation

The SDK exposes one public client package and keeps generated protocol code
internal. Standalone compilation and affected-target semantic lint pass.
All 34 selected repository-quality, owner OpenSpec, module-format, and DNS
consistency tests pass. No live API credentials or paid resources are used.

The provider worktree consumes an exact projection of the prepared SDK runtime
inputs. SHA-256 comparison of Go sources, generator configuration, and BUILD
files confirms identical inputs. Its real Terraform CLI acceptance suite passes,
including lifecycle, duplicate-name protection, remote deletion, management-key
revocation retry, malformed metadata, and multiple teams under one bootstrap
account. The reusable module is verified with two teams and three project keys.

SDK, related tools/schema configuration, and the reusable Terraform module are
the prerequisite PR. Provider and bootstrap integration remain in PR 249. The
dependent provider will be rebased onto trunk after this prerequisite merges;
this source change grants no merge or deployment authority.

The Options constructor is consumed by the provider and independent team-key
clients. The existing Terraform CLI lifecycle fixture now uses the role name
test_engineering/read. Before the transport fix, its first generated read failed
because the escaped slash became a path separator. With the fix, the complete
acceptance suite passes, including import, permission updates, and deletion.
The team-resources.json artifact records the encoded GET and DELETE paths.
