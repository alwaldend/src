# Validation

The retained isolated TLS JSON-RPC fixture passed through
`//infra/xcp_ng/cmd/xo_config:xo_config_test`. It covers missing/ambiguous
inventory, conflicts, creation and matching repetition, lost creation response,
wrong postconditions, API failures and credential exclusion. Bazel retains the
repeatable result in `bazel-testlogs/infra/xcp_ng/cmd/xo_config/xo_config_test/test.xml`.
The helper no longer includes the Ansible cache-device command.

The bootstrap UNC path now matches the AppRole-prefixed NAS share
//nas.alwaldend.com/src_infra_dev_runner_cache. NAS/XCP-ng formatting/provider
checks and the retained storage integration fixture pass. The creation-only
bootstrap still leaves retained SR/PBD settings unchanged; any existing mount
requires separately authorized reconciliation rather than automatic replacement.

XCP-ng provider validation and formatting passed. Source review confirms the
creation-only bootstrap preserves existing storage and delegates the repository
to the development resource set. The account now comes from NAS-managed Vault
publication. Plans do not reconcile retained PBD settings or rotate credentials.

The [coordinated receipt](../../../../dev/openspec/changes/mount-runner-cache/validation.md)
records quality, semantic lint, packaging and strict specification checks. No
SR, VM disk or host mount was created; live encrypted storage verification is
still required after separately authorized deployment.

The credential injector now reads shared/read/approle/src_infra_xcp_ng/nas using
its own AppRole. The cross-account Vault policy and binding are removed. The
development NAS project supplies cache filesystem/share access and parent traversal;
no live SR credential or mount was changed.
