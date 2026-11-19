## 1. Shared DNS roots

- [x] 1.1 Package the common Terraform configuration and both provider variants through a shared Bazel entry point.
- [x] 1.2 Adopt the entry point in Dev and Local Inference and document the owner inputs.

## 2. Validate

- [x] 2.1 Verify both packaged roots, including provider selection, and run offline checks.

Bazel built both entry points and the DNS documentation. The DNS declaration
and Terraform format checks passed. With HTTP initialization disabled, both
roots passed Terraform validation. Mocked plans against the previous and
shared root configurations confirmed the same 12 Dev and 3 Local Inference
resource addresses. Local Inference packaged only the Cloudflare provider.

## 3. Complete adoption

- [x] 3.1 Use shared packaging for every existing DNS-only root.
- [x] 3.2 Extract mixed DNS roots and move all declarations into their owners' `dns/` packages.
- [x] 3.3 Remove interim transfer code and document the required rollout order.
- [x] 3.4 Verify relocated inputs and builds offline.

Before the Work1 retirement, all 19 inventories were moved without content
changes. Offline Terraform validation covered 18 DNS roots and ten service
roots. Mocked DNS plans confirmed all 149 resource addresses.

## 4. Retire the 3x-ui Work1 node

- [x] 4.1 Remove the retired node's DNS declarations, Terraform configuration, and inventory entries.
- [x] 4.2 Delete the four Work1 records and verify the public and dc1 DNS answers.
- [x] 4.3 Transfer the remaining 12 3x-ui DNS bindings and verify a no-change plan.
- [ ] 4.4 Reconcile the retired node's service state after approval.
