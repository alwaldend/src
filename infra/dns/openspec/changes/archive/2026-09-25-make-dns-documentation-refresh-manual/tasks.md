## 1. Decouple documentation refresh

- [x] 1.1 Record failure cases and demonstrate the existing stale-page failure before changing validation.
- [x] 1.2 Remove ordinary freshness enforcement while keeping declaration and renderer validation; document explicit manual commands.
- [x] 1.3 Verify stale snapshots do not block or change during ordinary checks, explicit refresh/check works, and declaration errors still fail.
- [x] 1.4 Validate the owner specification, formatting, semantic lint, and repository quality before publication.

## Evidence

The stale-page probe failed before the change at
`TestGeneratedDNSPagesAreCurrent`. After the change, `//infra/dns:config_test`
passed with the same stale page and did not rewrite it; the original page was
restored afterward. The isolated CLI fixture verified `--check` failure,
explicit `--write`, successful subsequent `--check`, correct per-view content,
and rejection of duplicate declaration ownership. No repository declaration
page was regenerated. Fixtures and logs are under `out/split/dns-manual/`.

The pre-publication source checks passed all 32 selected tests, including
repository quality, DNS validation, Vault formatting, and both owner OpenSpec
workspaces. Semantic lint passed for the six affected owner packages. Delivery
receipts bind the final candidate and repeat the required publication gates.
