## Source acceptance

Terraform declares work1 with its own provider alias and shared node module.
Existing node resource addresses remain unchanged. The inventory and owned DNS
configuration use the same hostname and supplied addresses. Credentials remain
references to work1 entries in the existing Vault maps.

Architecture changes are preserved in the separate move-threexui-architecture
change and branch. Exact-candidate checks are recorded under ignored out/work1.
Subsequent explicitly authorized deployment evidence is recorded separately under
out/work1/deploy and does not change the offline validation scope.

The live Mullvad relay inventory no longer lists de-dus-wg-003. The replacement
de-fra-wg-001 is active and its endpoint and public key match the pinned data.
Only work1's selection changes; live verification receipts remain task-local.
