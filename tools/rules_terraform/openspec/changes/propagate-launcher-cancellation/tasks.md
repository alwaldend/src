## 1. Cancellation propagation

- [x] 1.1 Reproduce the held output pipe using the rental command, actual launcher, and signal-aware Terraform fixture before implementation.
- [x] 1.2 Forward graceful cancellation through both command boundaries and verify the fixture reaches cleanup and writes its receipt.
- [x] 1.3 Document the cancellation bounds and pass the launcher, rental lifecycle, and owning OpenSpec checks.
