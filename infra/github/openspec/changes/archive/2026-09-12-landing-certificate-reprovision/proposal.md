## Why

A project landing site can serve HTTP while its GitHub Pages TLS certificate is
stuck in a failed or stale state. GitHub Pages issues the custom-domain
certificate as part of its Pages build and does not expose a direct
"reprovision certificate" action. The only supported recovery is to remove and
re-add the custom-domain Pages configuration so GitHub schedules a new
certificate.

The current documented workflow only enables Pages for new sites and explicitly
forbids naming an existing served site in `bootstrap_projects`. That leaves no
sanctioned, reviewable recovery path when a live site's certificate needs
reissuing.

## What Changes

- Document and support a certificate reprovision procedure that removes a live
  site's custom-domain Pages configuration and then restores it, forcing GitHub
  to request a new certificate.
- Add a distinct `certificate_reprovision_projects` variable so this
  temporary-disabling intent is separate from first-creation
  `bootstrap_projects` and cannot be confused with it.
- Keep the guardrail that forbids using `bootstrap_projects` for an
  already-served site; the recovery path is explicit and differently named.
- Publish the built site content again only if the reprovision does not
  preserve it, and record how to verify the issued certificate.

## Capabilities

### New Capabilities

- None.

### Modified Capabilities

- `infra-github`: Add a requirement for the supported certificate reprovision
  procedure and the separate intent variable that gates it, and record the
  pinned published hostname that keeps a source-only rename from destroying a
  live landing repository.

## Impact

- `infra/github/tf/alwaldend_pages_repos.tf`: new variable and its effect on the
  `github_repository_environment.project_landing_pages`/Pages block `for_each`.
- `infra/github/tf/BUILD.bazel`: any registry-derived inputs if the procedure is
  parameterized there.
- `infra/github/tf/README.md` and `infra/github/README.md`: document the
  recovery workflow and its verification step.
- `projects/hugo_landing/skills/add-project-site/SKILL.md`: keep the rollout
  guidance consistent with the new, separately named variable.
