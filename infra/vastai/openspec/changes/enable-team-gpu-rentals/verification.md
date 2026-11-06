## Current status

The source uses a local GPU child module with an explicitly passed
`vastai.alwaldend` provider. Its key comes from a separate Vault reference.
Ordinary rentals and the paid lifecycle test have independent state selectors.
The rental source performs no context switching or credential creation. A
team catalog declares the full-access key through its shared team module and
stores the sensitive output in Vault. The imported alwaldend identity remains
in this root's existing state; its address reconciliation into the team module
is pending. AL supplies the issuer and both Vault provider credentials. The
separate credential root and backend declaration are removed.

## Offline evidence

- The pinned Terraform CLI validates the actual rental root and child module
  with a loopback HTTP backend and dummy credentials.
- Ten command scenarios pass: success, absent backend state, backend failure,
  malformed instance inventory, personal account, partial creation failure,
  failed release, existing rental state, a different team, and cancellation
  through the actual packaged Terraform launcher.
- The command receipt records only non-sensitive lifecycle facts. Fixtures
  verify one creation attempt followed immediately by release, including on
  partial failure. Separate cases verify that preflight failures do not mutate.
- Terraform formatting and owner OpenSpec validation pass. Packaged launchers
  include the child module and the exact `tf=rentals`, `tf=rental_test`, and
  `vastai=rentals` selectors. The owning apply no longer requires a saved plan.

Final repository quality, semantic lint, and local review are delivery gates.
Their candidate identities and publication status belong to delivery receipts.

## Live credential observations, 2026-11-06

Read-only verification fetched the two configured Vault entries without
printing or retaining their values or full account responses:

- `src_infra_vastai/vastai` selects the personal account.
- `src_projects_vastai_terraform_provider/vastai` selects team
  `test_vastai_provider`, which has no credits.
- The proposed `src_infra_vastai/vastai/alwaldend` entry is absent.
  Listing the component namespace found only its backend directory and existing
  personal credential. A broader Vault root listing was unavailable; this does
  not establish that no replacement exists elsewhere.

The previous reconciliation record identifies the old manual alwaldend team
key as revoked. No replacement was found among the configured Vault references.
The user requires a full-access team key in Vault; a different team's key or a
personal-account key does not satisfy that requirement. The live rental test
has not run and no GPU has been rented.

An initial context-bootstrap attempt generated key ID 36342752 and failed
activation. A reviewed owning Terraform destroy revoked that exact key and
completed with one deletion and no other resource changes. Its implementation
was removed after the user selected direct-key child-module wiring.

The subsequent user-authorized credential root planned two additions: one API
key and its Vault entry, with no team or credit changes. After reviewing that
interactive plan, applying with the existing team's ID failed at
`POST /auth/apikeys/` with HTTP 401. No API-key ID was returned and the Vault
resource was not created. That failure was specific to the direct permission
override request; the issuer has key-creation permissions.

The corrected full-access request omits the permissions field (an empty JSON
object first received HTTP 400). Using the same stored issuer, the next reviewed
apply created key 36348922 and retained it in state. Selecting the alwaldend
owner context on that generated key returned HTTP 400. The Vault resource did
not run. A reviewed owning destroy revoked exactly key 36348922: zero additions,
zero changes, one destruction. The shared issuer was never switched. The actual
remaining defect is team binding, not absence of the existing issuer. No
replacement Vault entry or GPU rental was created.

## Session ergonomics

- `vastai-credential-context-assumption` (live): introducing credential
  generation before verifying the stored keys' account contexts caused an
  unnecessary provider implementation. Read the actual Vault references and
  verify selected accounts before choosing a credential workflow.
- `vastai-empty-state-preflight` (fixture-tested): real Terraform returns
  exit 1 for an absent HTTP-backend state. The command now accepts only that
  explicit diagnostic as empty and preserves other backend failures.
- `bazel-wrapper-default-arguments` (configured): directly invoking an AL
  wrapper executable omitted Bazel's fixed arguments. Use its owning run target.
- `formatter-output-volume` (configured): the first formatter emitted a large
  path inventory. Later runs retain diagnostics under ignored task scratch and
  report bounded results.
- `terraform-launcher-cancellation` (fixture-tested): a SIGTERM to the rental
  process reproduced an eight-second held-pipe failure through the real launcher.
  Graceful forwarding at both boundaries now reaches the child's interruption
  checkpoint, performs release, and writes the sanitized cancellation receipt.
- `vastai-team-key-binding` (live): the existing issuer successfully creates
  keys, but selecting the existing team on a new default full-access key returns
  HTTP 400. General issuer-inadequacy claims were incorrect; preserve the exact
  failing operation and retain partial keys for explicit IaC cleanup.
- `vastai-credential-module-duplication` (fixture-tested): a parallel credential
  root repeated the shared module's key lifecycle and issuer input. The shared
  module now owns that lifecycle through imported identity. Provider-default authentication, legacy issuer
  removal, and raw-key-to-module migration pass without key rotation. No
  unchanged failed live request was retried for this refactor.

These are task-local observations, not changes to shared policy or skills.

## Superseded separate-root module revision

The shared fixes PR replaces the existing-team selector with import into the
module's singleton team resource. The caller now derives its team ID from
resource state and retains no temporary migration block. The unused raw-key
move was removed: both failed keys were already revoked and their credential
state contains no managed key to move. Cleanup targets only the key and Vault
entry so imported team identity is preserved.

Alwaldend's live team identity remains in the bootstrap state. The source
assigns it to the credential module, but ownership transfer is pending and must
precede any credential apply. No live state migration, key issuance, or rental
was run for this source revision. The previously observed two-addition key
plan describes the former selector-based candidate, not acceptance of this
unapplied ownership migration.

## Prior prerequisite merge and rebase

PR #267 merged into master at f6be2da53330726fb4e9ad1ed11dd73cffa7681b.
Its shared provider, module, infrastructure-skill, and completed bootstrap
migration cleanup source exactly matches the preserved integration inputs.
The dependent rental candidate is rebased onto that commit; those prerequisite
changes belong to master rather than this PR. The existing alwaldend ownership
transfer and live binding acceptance remain pending.

## Catalog consolidation

The extra credentials root was unnecessary because the existing team catalog
and shared team module already own team/key composition. Alwaldend is restored
to that catalog with its rental key. The existing project-secret address and
Vault identity remain unchanged; the rental secret uses the component's AppRole
through an aliased Vault provider. No broader policy or cross-role secret access
was added. The retired credential state has no key after the recorded cleanup.

The real-root Terraform fixture first failed because alwaldend and the component
Vault route were absent. It then verifies both mocked Vault identities and
rejects duplicate destinations and unsupported identity names. This evidence
covers offline configuration behavior, not live key binding. The remaining
state transition is an address move within the existing bootstrap state, rather
than the previously proposed transfer to a separate credential backend. No live
migration, key creation, or paid rental was performed for this revision.

The repeated separate-root design was a lifecycle-boundary error: calling the
shared team module from another root still imposed an unnecessary backend and
ownership transfer. The corrected boundary is one catalog root for teams, keys
and secret destinations, with provider aliases handling Vault identities.
