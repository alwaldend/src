---
title: Vast.ai rental test
description: Rent the cheapest on-demand GPU in alwaldend and release it
---

This flat Terraform root owns one test contract, independently of account and
team IaC. The provider reads `VAST_API_KEY`, injected from the generated team
key at `secrets/alwaldend.com/vault1/approles/src_infra_vastai/vastai/alwaldend/rentals`.
The [alwaldend root](../alwaldend/README.md) owns that key. The backend retains
the existing `tf_backend/rental_test` path.

```sh
bazel run //infra/vastai/rentals_test:tf.plan
bazel run //infra/vastai/rentals_test:tf.rent
bazel run //infra/vastai/rentals_test:tf.release
bazel run //infra/vastai/rentals_test:tf.state -- list
bazel run //infra/vastai/alwaldend -- show instances
```

Review both interactive applies. Start release immediately after rent returns,
including after a failed creation apply, and approve destruction promptly.
There is no boot wait or workload. Confirm empty Terraform state and absence of
the contract in the CLI inventory. Release destroys the contract and local disk;
stopping a contract alone retains storage charges.

The search takes one rentable, unoccupied on-demand GPU offer ordered by total
hourly price. It includes the 3.5 GB disk allocation used by `ubuntu:24.04` and
rejects prices above USD 1 per hour. Network charges are separate. There are no
model or verification filters and no creation retries. An existing contract
retains its original offer; `selected_offer` describes the current marketplace.
Destroy mode releases the recorded contract without searching for another offer.

Building or testing the package does not rent compute. The offline catalog test
validates this root with the packaged provider. Provider acceptance tests cover
the Terraform rental lifecycle against a local API fixture.
