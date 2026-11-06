---
title: Vast.ai GPU rentals
description: Cheapest on-demand GPU rentals in the funded alwaldend team
---

These targets pass the aliased `vastai.alwaldend` provider to the local GPU
child module. Its dedicated team key is read from field `api_key` at
`secrets/alwaldend.com/vault1/approles/src_infra_vastai/vastai/alwaldend` and
injected as `TF_VAR_alwaldend_api_key`. No credential context is switched.
The team key stored in Vault must have full team access. The existing
[team catalog](../tf/README.md) creates it through the shared team module and
owns Vault storage. Its generated key could not select alwaldend (HTTP 400), so rental
access remains blocked pending a working team binding and key verification.
The personal bootstrap key remains in its separate existing Vault entry.

Ordinary rentals and the explicit live test share the module but use different Vault-backed states:
`tf_backend/rentals` and `tf_backend/rental_test`.

## Rent and release

```sh
bazel run //infra/vastai/rentals:tf.plan
bazel run //infra/vastai/rentals:tf.apply
bazel run //infra/vastai/rentals:tf.release
```

Review each interactive apply and enter `yes` only for the intended contract.
No saved plan or automatic approval is used. Release destroys the contract,
including its local storage; stopping an instance alone retains storage charges.

Selection searches rentable, unoccupied on-demand GPU offers, orders by
`dph_total`, and takes the cheapest result. It has no GPU-model or verification
filter. The default disk allocation is 1 GiB and is included in the offer-price
calculation through `allocated_storage`. Default image is `ubuntu:24.04` and
the configurable total hourly price ceiling is USD 1. Network traffic charges
are separate. Offers can disappear before apply; creation is not retried.

Use identical variable overrides for plan and apply, such as `-var=disk_gb=20`
or `-var=image=...`. An existing rental retains its original offer even when
marketplace results change. Other creation-setting changes require replacement.
The `selected_offer` and `rental_hourly_price` outputs describe the current
marketplace result, not a refreshed billing record for an existing contract.

Apply acknowledges the contract without waiting for boot or SSH. A later
`tf.refresh` observes connection fields. The lifecycle test needs no SSH key;
interactive workloads require account-side SSH setup before connecting.

## Explicit paid lifecycle test

```sh
mkdir -p out/vastai-rental
bazel run //infra/vastai/rentals:live_test -- \
  --receipt "$PWD/out/vastai-rental/live-receipt.json"
```

The receipt path must be new. The command verifies the funded team context,
refuses nonempty test state, and records existing contract IDs without exposing
account details. It then runs an interactive creation apply. Immediately after
that command returns, it starts an interactive destruction apply, even if
creation failed. Review and approve that single-contract deletion promptly;
there is no boot wait, SSH check, or workload execution between the applies.
This is a billable operation and is never part of `bazel test //...`.

After release, the command checks empty test state, confirms that the contract
is absent from the live instance inventory, and checks that earlier unrelated
contracts remain. The JSON receipt contains only timestamps, team name,
contract ID, selected hourly price, and lifecycle results. An ambiguous creation
or cleanup failure fails the test and requires account inspection before retry.

SIGINT or SIGTERM received by the command during creation cancels creation and
attempts cleanup with an independent five-minute deadline. The launcher forwards
an interrupt to Terraform and allows ten seconds for it to exit; the caller's
twenty-second fallback bounds its wait and output pipes. Parent-wrapper
termination, a hard kill, failed authentication,
declined confirmation, or network/backend outage can prevent cleanup. Inspect
the receipt and recover through this test state's owning targets:

```sh
bazel run //infra/vastai/rentals:test_tf.state -- list
bazel run //infra/vastai/rentals:test_tf.release
```

These recovery commands use the test state, while `tf.release` uses ordinary
rental state. Neither root owns the team or its credentials. Offline command
fixtures cover success, partial creation failure, cleanup failure, existing
state, and incorrect team selection; provider acceptance covers real Terraform
lifecycles against loopback API fixtures.
