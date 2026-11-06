## Current design

The account root `infra/vastai/tf` uses the personal account credential to own
team identities and account environment settings. Its desired team set contains
only alwaldend. It creates no project API keys or rental contracts.

The separate root `infra/vastai/alwaldend` uses an operator-provided full-access
team key from Vault as its default provider credential. It creates full-access
project keys in the selected team and stores them under the component namespace.
It also applies the shared account-variable map to the team. The seed key is
operator-owned and is not replaced by Terraform.

`infra/vastai/config.json` owns the environment map and declared teams/key names.
Account and team roots read the same environment section. Teams do not inherit
account variables. Current values are ENABLE_HTTPS=true, WEB_ENABLE_HTTPS=true,
and WEB_ENABLE_AUTH=true.

The account, team and rental test backends are independent:
`tf_backend/tf`, `tf_backend/alwaldend`, and `tf_backend/rental_test`.
No resource or credential is copied between states. The rental test state was
verified empty after live acceptance, so flattening requires no address migration.

## Authentication evidence

The stored account API key has full personal-account rights. Generated API keys
cannot switch contexts: the live request returned HTTP 400 with session-only
rejection categories. A direct custom-permission override request returned HTTP
401; no key was created. These failed experiments do not establish an alternative
bootstrap contract and their speculative implementation was removed.

The pinned CLI exposes key creation without a team selector and instructs
operators to select an existing team in the console to create its initial key.
The user provided that key through the component Vault command. CLI inspection
confirmed alwaldend account 846045, full team permission groups, positive credit,
and an initially empty instance inventory. Ordinary IaC uses provider credentials
without switching context or supplying a resource-level issuer.

## Retirement and state reconciliation

Destroy the former live test root's key, marker and role before revoking its
credential. Then remove other teams and their generated credentials from desired
account configuration. Preserve alwaldend, its funded identity and operator keys.
Imported identities without recovered management keys cannot currently be
destroyed by the provider; retain their state and report incomplete deletion.

Reconcile the alwaldend identity and account variable addresses in the existing
account backend. Remove temporary migration blocks after state verification.
The separate team root owns newly created resources only.

## Rental lifecycle

The flat `infra/vastai/rentals_test` root uses the generated rental key through
the default provider's VAST_API_KEY environment variable. It has no child module,
credential variable, duplicate rental backend or lifecycle command runner.
Search rentable, unoccupied on-demand GPU offers, order by total hourly price,
and select one result. Include disk allocation in the comparison and enforce a
USD 1/hour price ceiling. The image is ubuntu:24.04; live acceptance
showed that a 1 GB request becomes 3.5 GB, so the default and price calculation
now both use 3.5 GB. Preserve an existing contract's original offer ID.

Separate `tf.rent` and `tf.release` commands review interactive creation and
destruction plans. The operator starts release immediately after rent returns,
including on creation failure, without waiting for boot or SSH. Destroy mode
uses the recorded contract without marketplace discovery. State listing and the
team CLI verify cleanup. The recorded live acceptance receipts remain evidence
of the preceding implementation; this simplification adds no paid rental.

## Verification

The actual account and team configurations run in separate isolated Terraform
working directories with mock providers. Their assertions cover the declared
team, shared variable map, default full-access key and separate Vault destination.
Invalid key names fail before provisioning. The same fixture initializes and
validates the actual flat rental root with the packaged provider. Live acceptance independently verifies
credentials, settings, rental creation/release, CLI inventory and owning states.
