---
name: add-project-site
description: >-
  Add or update a projects/* landing site using the shared Hugo page,
  project registry, GitHub Pages Terraform, and project DNS declaration.
  Use for project-site onboarding and explicitly requested deployment in this
  repository; do not use for unrelated website design or apex-site changes.
---

# Add a project landing site

Read the [Hugo Landing README](../../README.md) for the supported build and
rollout entry points. Inspect the project's README and existing landing before
editing; the published page must preserve its actual purpose, status, and
limitations. Improve source content where requested without inventing features
or presenting unfinished work as released.

## Wire the site

- `projects/projects.bzl` owns `PROJECTS`, the dedicated-site registry. The
  `alwaldend.com` project uses the existing apex website and stays outside this
  list. Derive coverage from the current registry rather than copying a count.
- Ordinary projects own `projects/<project>/landing:site`, using
  `al_hugo_landing_site` and their public README/docs target. Keep content in
  the owning project and shared presentation in `hugo_landing`.
- Nested `rules_*` modules use root-workspace assembly at
  `//projects/hugo_landing/landing:<project>` and external labels such as
  `@rules_example//:docs`. An exported root README is sufficient when the
  module has no docs target. Ensure the header's Docs URL is actually published
  by `//projects:docs`, or point it to the source README. Do not make a standalone module depend on the root
  repository just to render its landing page.
- Add the registry entry, landing target, project-owned `dnsconfig.json`, and
  its BUILD export/visibility together. Follow a neighboring project for the
  unproxied global CNAME directly to `alwaldend.github.io.`. Project identifiers retain
  underscores; hostnames and `<project>-landing` repository names use hyphens.
- `//projects:landing_sites` and `//projects:landings` aggregate the registry;
  `infra/github/tf` generates repository declarations from it; `infra/dns`
  collects each project's DNS target. Inspect these consumers after changing
  membership. Do not hand-edit generated mappings or duplicate apex records.
- Add the site's clickable link to `projects/README.md`. Include a useful
  description and introduction in the project README. Use the shared Docsy
  header; do not duplicate navigation as buttons or add generated-source labels.

## Check the candidate

Use `repo-bazel` and `bazel-agent` for the owning build/tests and formatter.
Build the affected rendered site, check `//projects:landings` and
`//infra/dns:config_test`, and inspect representative HTML from the candidate.
Check its project identity, useful description, documentation/source links,
relative links and images, and preserved code examples. A successful build
alone does not establish a usable or deployed page.

## Deploy when authorized

Adding a site authorizes implementation and offline checks. Live Terraform,
GitHub publication, and DNS mutations require the user's requested operation
and scope; retain authorization already given instead of requesting it again.
Use `repo-terraform` and `repo-secrets` for their owning procedures. Follow the
README rollout entry points, including `//projects:deploy_landings` with an
absolute task-owned `--scratch` directory. Keep a single-project request
scoped to that project when choosing deployment arguments.
Use `--project <project_id>` on the aggregate publisher for a single site.

New repositories need a `pages` branch before Terraform enables Pages. The
`bootstrap_projects` variable defaults to an empty set. For first creation,
plan with only the new projects in that set, inspect the plan, and apply only
changes within the authorized scope. Never include an existing served site
just to get bootstrap past a failure. Publish the rendered pages branches,
then plan/apply with the default empty set to enable Pages for those sites.

Preview DNS changes using the affected record names and the global provider;
apply the same filtered scope only after checking the preview. Do not include
unrelated infrastructure or apex changes. If the plan or preview expands the
scope, resolve the configuration or selection before executing it.
The public virtual zone is `alwaldend.com!global`: use
`--domains 'alwaldend.com!global' --providers global` with the preview/deploy
targets. The bare `alwaldend.com` selector matches zero configured zones;
a zero-zone preview does not establish that DNS is current.

Verify the pushed revision, Pages configuration/build state, DNS, and HTTPS
page content. Report observed successes separately from pending propagation,
certificate issuance, or failed operations. A Terraform apply, Git push, or
DNSControl exit code alone is not evidence that the public site is serving.
