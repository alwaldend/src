---
name: repo-hugo
description: >-
  Work with Hugo in this repository: the pinned toolchain and al_hugo_* rules,
  the checked-in theme lock and Hugo module mounts, and the shared landing-site
  shell. Use for Hugo site or theme changes, including the apex site, and for
  project landing-site onboarding and explicitly requested deployment; do not
  use for unrelated website design.
---

# Hugo sites in this repository

Read the owning READMEs before editing: [Hugo Landing](../../../hugo_landing/README.md) for the
shared landing shell and rollout entry points, [Rules Hugo](../../../rules_hugo/README.md)
for the toolchain and build rules, and [Alwaldend.com](../../../alwaldend.com/README.md)
for the apex site. The [Hugo miscellaneous page](../../../alwaldend.com/content/docs/misc/hugo.md)
records the environment-allowlist constraint that the site build depends on.

Keep one implementation of each shared behavior. The apex site and project
landings share one toolchain, one theme lock, and one presentation shell; a
site-local copy of a shared style, layout, configuration fragment, or build
rule is a defect. Extend the shared owner, or parameterize it, instead of
duplicating its content. Fix every reference at the owning source rather than
repeating a fact in another file.

## Build model

- Hugo is pinned as a Bazel toolchain. `al_hugo_extension` in
  `projects/rules_hugo` downloads the requested release and registers a
  `@rules_hugo//pkg/bzl:toolchain_type` toolchain per platform. The site
  archive stays in the target configuration while the Hugo binary stays in the
  execution configuration.
- `projects/alwaldend.com/include.MODULE.bazel` declares the toolchain version,
  the theme lock, and the `use_repo`/`register_toolchains` calls. Treat it as
  the owning declaration for this workspace and follow it instead of inventing
  a new Hugo dependency.
- `al_hugo_site` pairs a site source archive with PostCSS and build tools,
  `al_hugo_worker` renders it through the persistent `@rules_hugo//cmd/hugo_worker`,
  and `al_hugo_binary` exposes a runnable Hugo command for a site.
- Sites are ordinary targets: `//projects/alwaldend.com:site` for the apex site
  and `//projects/<project>/landing:site` for each project landing. The apex
  site additionally builds under `--config=release`, which sets the public
  `--baseURL`, and `//projects:deploy_landings` publishes rendered sites.
- Build flags that matter: `--panicOnWarning` fails the build on any Hugo
  warning, and `--printPathWarnings` reports unresolved paths. The apex site
  always builds expired and future content, and only its local preview adds
  `--buildDrafts`, so a `draft: true` page appears in preview but never in a
  release build. Landing sites pass neither flag and render only published
  content. Do not silence these flags to get a build through; fix the source.
- Local preview: `bazel_agent bazel run //projects/alwaldend.com:site_serve`
  serves the built site on `127.0.0.1:1313`. Landing sites have no preview
  target; inspect their rendered output directly.

## Use a theme

- A theme is consumed as a Hugo module named by its repository path, for
  example `theme = ["github.com/google/docsy"]` in the site configuration.
  The generated landing configuration and the apex `hugo.toml` both select
  Docsy this way.
- The theme, Bootstrap, and Font-Awesome sources are pinned by the checked-in
  lock at `projects/alwaldend.com/hugo_lock.json`. Each entry carries the
  module `path`, archive `url`, optional `root` and `strip_prefix`, and the
  `integrity` digest. The extension turns every entry into an archive
  repository and a `themes` filegroup; `include.MODULE.bazel` exposes those
  repositories through `use_repo`.
- `[[module.mounts]]` in the site configuration maps those downloaded trees
  into the Hugo module layout: the site's own `assets`, Bootstrap `scss`, `js`,
  and `dist`, and Font-Awesome `scss` and `webfonts`. A theme is only usable
  when every path its stylesheets and scripts reference is mounted.
- Project styles are picked up by Hugo's theme extension points, not by
  overriding theme files: `assets/scss/_variables_project.scss` for theme
  variables and `assets/scss/_styles_project.scss` for added rules. The shared
  landing shell keeps its styles in
  `projects/hugo_landing/assets/scss/_styles_project.scss`; the apex site keeps
  its own in `projects/alwaldend.com/assets/scss/`.
- Changing a theme version means updating the lock entry and the mount or
  configuration that depends on it together, then rebuilding the affected
  sites. Never edit a downloaded theme tree; it is generated from the lock.

## Add a project landing site

Read the [Hugo Landing README](../../../hugo_landing/README.md) for the supported build and
rollout entry points. Inspect the project's README and existing landing before
editing; the published page must preserve its actual purpose, status, and
limitations. Improve source content where requested without inventing features
or presenting unfinished work as released.

### Wire the site

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

### Check the candidate

Use `repo-bazel` and `bazel-agent` for the owning build/tests and formatter.
Build the affected rendered site, check `//projects:landings` and
`//infra/dns:config_test`, and inspect representative HTML from the candidate.
Check its project identity, useful description, documentation/source links,
relative links and images, and preserved code examples. A successful build
alone does not establish a usable or deployed page.

### Deploy when authorized

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
