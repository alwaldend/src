## 1. Implementation

- [x] 1.1 Move canonical DNS and firewall configuration, wire separate backends and dedicated credentials, and verify formatting, provider schema and DNS ownership; document state transfer and acceptance.

- [x] 1.2 Replace RouterOS credential variables with native provider environment injection through the shared helper; verify provider schema validation and inspect credential paths.

- [x] 1.3 Add and render a network diagram showing LAN and Wi-Fi actors, router and printer DNS names, the printer-only firewall exception and a bidirectional printer connection; inspect both theme variants and the documentation asset wiring.

## 2. Verification

- [x] 2.1 Validate the exact candidate with repository quality, affected semantic lint and strict OpenSpec checks, recording results and live-operation limits.

- [x] 2.2 Revalidate and locally review the updated candidate, including the provider-environment skill guidance.

## 3. Authorized deployment

- [x] 3.1 Check live source bindings and exact provider IDs, provision the dedicated account and AppRole, and preserve existing printer objects during binding transfer.
- [x] 3.2 Verify no-change plans, rule order and network reachability; remove temporary migration tooling and retain a sanitized deployment receipt.
