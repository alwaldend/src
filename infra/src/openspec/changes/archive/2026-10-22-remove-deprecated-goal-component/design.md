## Context

The 2026-09-08 migration moved all nine maintained goal records into owner
OpenSpec workspaces and left `projects/goal` as deprecated compatibility code
with a disabled skill. Since then the pinned OpenSpec CLI has been the
documented channel, and `AGENTS.md` already told agents not to create goal
records.

The component is not only Go source. `projects/projects.bzl` lists it, which
drives landing-site rendering, per-project DNS-config packaging, and GitHub
Pages repository generation. `projects/goal` also declared its own `dnsconfig`
CNAME and `docs/`, and the rendered zone snapshot and the site test both
reference its documentation pages.

## Goals / Non-Goals

**Goals:**

- Remove the goal tool, store, skill, docs, landing site, and DNS record.
- Keep the generated project lists, zone snapshot, and site checks consistent.
- Preserve the historical record without reviving the tool.

**Non-Goals:**

- Reinterpret, revalidate, or rewrite migrated acceptance history.
- Change the OpenSpec workflow or the owner-local workspace layout.
- Deploy or reconcile live DNS as part of this change.

## Decisions

The `goal` entry is removed from `PROJECTS` because that list is the single
source for the landing site, the project DNS config set, and the Terraform
Pages map. Leaving it would keep generating DNS and Pages configuration for a
project with no source.

The committed zone snapshot drops its `goal` CNAME line. The zone files are
rendered during deployment; removing the matching `dnsconfig.json` without
updating the snapshot would leave the checked-in record inconsistent with the
configuration. No provider call is made, and the record is not deleted from
live DNS here.

`projects/goal/openspec/specs/project-goal/spec.md` and the `projects_goal`
validation entry are removed, and the `repository` baseline requirement changes
from "deprecated and disabled" to "removed". The `project-goal` capability has
no component left to describe.

The site test's goal-specific page and image regression block is removed with
the pages it asserted.

The removed bytes are not copied or restated in this change. They remain in git
history at the removal commit's parent, so duplicating them or their digests
would add no guarantee.

## Risks / Trade-offs

- The `goal.alwaldend.com` DNS record is removed from configuration but still
  resolves until a deploy runs. Reconciliation is a separate authorized
  infrastructure operation.
- The `goal` GitHub Pages repository and its published site are left in place;
  removing them is a forge operation outside this change.
- An external consumer of the goal CLI would break. The deprecation notice
  since 2026-09-08 is the only migration affordance, and no in-repository
  consumer remains.
