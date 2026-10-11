## Context

Five component-owned dependency sets serve repository Python tools, Ansible,
Ykman, Autoscroll and firmware tooling. Their manifests and native locks belong
beside their owning BUILD files. The shared root pyproject.toml owns formatter
settings only. Firmware source remains pinned externally; its npm extension
identity preserves lazy loading of the upstream Node lock.

## Decisions

**Proceed:** use rules_python 2.4.2 and its matching Gazelle plugin. The public
`@rules_python//python/uv:lock.bzl` API selects native `uv lock` for `out =
"uv.lock"`. Each component loads it directly and owns `uv_lock.update`. The root
alias points to the repository tools lock. No shared wrapper or custom freshness
test is needed. Normal builds consume checked-in locks; explicit resolution
actions require network access and remain sandboxed.

`pip.parse(uv_lock = ...)` consumes each native lock directly. Gazelle uses the
same tools/py lock through its `lockfiles` parameter. Requirements exports and
maintained pip locks are removed. The upstream hub's generated requirements.bzl
labels remain its public API, rather than a separately maintained lock format.

Select Python 3.14.4, the newest stable runtime bundled in rules_python 2.4.2.
Its bundled Python 3.15 version is still a release candidate. The user's explicit
instruction to use the rules' available runtimes excludes a custom archive
override. uv remains pinned to 0.13.0. Both executables are Bazel-managed; no host
interpreter or global configuration is changed.

Preserve compatible existing dependency pins while adopting universal locks.
Migration-only constraints seed the native resolution and are removed from the
final manifests. Native locks may include additional Python/platform branches;
the repository toolchain selects Python 3.14. Upgrades needed by that runtime
are inspected and verified through real consumers. Autoscroll retains Python
3.9 support while selecting newer PySide6 on Python 3.14.

The firmware tools manifest and lock now belong to projects/infinitime. Its
source patch removes the upstream legacy compile rule; it does not copy a new
locking workflow into the external repository. The existing hub name, source
commit and SDK/compiler pins remain intact. Its local Python lock now uses the
standard pip extension, retaining upstream repository mappings and platform
defaults. The obsolete re-export shim is removed; only npm retains an isolated
identity for its externally owned lock.
Historical patch.patch is inactive and is not a maintained dependency workflow.

The K3s role declares uv and python3 packages, creates its existing virtual
environment with uv, and installs kubernetes through uv pip. Its changed_when
contract follows sibling roles. Only packaging and a source-derived fixture in
ignored local scratch are authorized for validation; OS package installation,
K3s installation and Kubernetes operations are excluded.

## Validation and delivery

Regenerate all five native locks through their owning update targets, inspect
versions and markers, and compare generated outputs byte for byte. Refresh
Gazelle's manifest and Bzlmod metadata through their owning workflows. Validate
Python tool CLIs, the recommender E2E, Autoscroll packaging, firmware tooling,
Python skill configuration and repository quality/semantic lint. Gazelle now
records Autoscroll's existing import root and generates its missing library/CLI
dependency edges. Its help formatter must handle zero-argument actions before
requiring an argument type; the CLI help check exposed that existing defect.
The generated
lock bytes are excluded from generic formatting in .gitattributes.

The source-derived role fixture using Python 3.14.4 reported changed=2, failed=0 on its
first run and changed=0, failed=0 on its second. Its executable paths and logs
are task-local under out/migrate-python-uv/k3s-e2e-314. No K3s server was installed.

Stable diagnostic `cached-python-absolute-symlink` identified an interpreter link
in the shared Bazel contents cache pointing into an unavailable output base.
Ignored task-local user.bazelrc disables contents reuse and changes the watched
RULES_PYTHON_REPO_DEBUG value to refetch Python. Shared cache and host
configuration were not changed. Diagnostic `firmware-pip-package-metadata`
identified a missing repository mapping in the re-exported pip extension after
the upgrade; using the standard extension removes that obsolete shim.
Diagnostic `uv-run-runfiles-output` records that
the upstream .run target writes in runfiles; .update is the maintained source
update workflow.

Git and repo_delivery receipts own final candidate identity, validation and
publication. Delivery uses required quality/lint gates, guarded synchronization,
an exact-lease push and the existing pull request. Rollback is a revert of the
task-owned commit, with no deployed state to reconcile.

Version inspection for Python 3.14 on Linux preserves all 76 root tool pins,
17 Ykman pins and 11 firmware pins. Ansible retains its package versions and
excludes ruamel-yaml-clib through its Python marker. Autoscroll selects PySide6
6.12.0 and its updated Qt closure; the old 6.10.0 branch remains for older
Python. Every downloaded lock artifact carries a SHA-256 digest, and no
migration-only constraints remain in the native locks. Ansible, ansible-lint
and Ykman version commands pass. Autoscroll help and an offscreen Qt SVG widget
check pass on Python 3.14.4/PySide6 6.12.0. The firmware DFU consumer builds. All 31 selected tests pass, including root
repository quality, Gazelle manifest integrity, the recommender E2E, Autoscroll
and firmware build tests, and Python skill offline configuration validation.
Final semantic lint and exact-candidate gates are recorded by delivery receipts.

## Session ergonomics

The native API inspection exposed two avoidable assumptions: .run writes in
runfiles, and constraints in standalone uv.toml did not seed this project lock.
Use the documented .update target, and seed migration pins through temporary
project configuration over each package's declared Python support. Both were
verified through generated artifacts; no migration constraints remain in source.
A whole-repository Gazelle run encountered unrelated Go discovery before failing
on Autoscroll's missing import-root declaration. After correcting that owner
configuration, scope generation to the changed component. Evidence is bounded
in the native lock, CLI, generation and delivery logs under task-local out/.
