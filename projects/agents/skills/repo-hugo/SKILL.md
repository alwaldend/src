---
name: repo-hugo
description: >-
  Work with Hugo in this repository: the pinned toolchain and al_hugo_* rules,
  the checked-in theme lock and Hugo module mounts, and the main site's shared
  presentation shell. Use for Hugo site or theme changes, including the apex
  site and the project landing section at /projects, and for explicitly
  requested deployment; do not use for unrelated website design.
---

# Hugo sites in this repository

Read the owning READMEs before editing: [Rules Hugo](../../../../tools/rules_hugo/README.md)
for the toolchain and build rules, [Alwaldend.com](../../../alwaldend.com/README.md)
for the apex site and its packaged content, and [Projects](../../../README.md) for the
project registry. The [Hugo miscellaneous page](../../../alwaldend.com/content/docs/misc/hugo.md)
records the environment-allowlist constraint that the site build depends on.

Keep one implementation of each shared behavior. There is exactly one Hugo
site, one toolchain, one theme lock, and one presentation shell; a copy of a
shared style, layout, configuration fragment, or build rule is a defect.
Extend the shared owner, or parameterize it, instead of duplicating its
content. Fix every reference at the owning source rather than repeating a fact
in another file.

## Build model

- Hugo is pinned as a Bazel toolchain. `al_hugo_extension` in
  `tools/rules_hugo` downloads the requested release and registers a
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
- There is one site target: `//projects/alwaldend.com:site`. It builds under
  `--config=release`, which sets the public `--baseURL`, and
  `//projects/alwaldend.com:deploy` publishes the rendered site to the apex
  Pages repository.
- Build flags that matter: `--panicOnWarning` fails the build on any Hugo
  warning, and `--printPathWarnings` reports unresolved paths. The site always
  builds expired and future content, and only its local preview adds
  `--buildDrafts`, so a `draft: true` page appears in preview but never in a
  release build. Because every project landing is part of this one build,
  `--panicOnWarning` means a landing content error fails the whole site and
  the blog deployment. Keep landing content small and validate with the site
  build before publishing. Do not silence these flags to get a build through;
  fix the source.
- Local preview: `bazel_agent bazel run //projects/alwaldend.com:site_serve`
  serves the built site on `127.0.0.1:1313`.

## Use a theme

- A theme is consumed as a Hugo module named by its repository path, for
  example `theme = ["github.com/google/docsy"]` in the site configuration.
  `hugo.toml` selects Docsy this way.
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
  variables and `assets/scss/_styles_project.scss` for added rules. The site
  keeps the shared canvas and accent in
  `projects/alwaldend.com/assets/scss/_shared_theme.scss` and
  `_shared_accent.scss`, and its own rules in `_styles_project.scss`.
- Changing a theme version means updating the lock entry and the mount or
  configuration that depends on it together, then rebuilding the site. Never
  edit a downloaded theme tree; it is generated from the lock.

## Add a project landing page

The main site publishes each registered project's visitor-facing page at
`/projects/<name>/`, and its repository reference documentation at
`/docs/projects/<name>/`. There is no per-project Hugo site, Pages repository,
DNS record, or certificate. Read the project README and any existing landing
before editing; the published page must preserve the project's actual purpose,
status, and limitations. Improve source content where requested without
inventing features or presenting unfinished work as released.

### Wire the landing page

- `projects/projects.bzl` owns `PROJECTS`, the landing registry. The
  `alwaldend.com` project is served by the apex site itself and stays outside
  the list. Derive coverage from the current registry rather than copying a
  count.
- The project owns its landing content in `projects/<project>/site/content/`
  as `_index.md` plus a `BUILD.bazel` that exports it with
  `docs_filegroup(name = "docs", prefix = "content/projects/<project>")`.
  That directory carries content only: no layouts, styles, or site
  configuration.
- `projects/alwaldend.com/content/projects/BUILD.bazel` packages those
  directories into `content/projects/<name>/`, deriving membership from
  `PROJECTS`, so a registered project is declared once. A registered project
  without content produces no page and no build failure.
- The apex site renders a landing through the front matter contract
  `layout: landing` with `title`, `linkTitle`, and `description`. The
  `projects/landing.html` layout and the `.landing-content` styles in
  `projects/alwaldend.com/assets/scss/_styles_project.scss` own its
  presentation. Never add a project-local layout or stylesheet.
- Landing front matter participates in the site's `statuses`, `languages`,
  and `tags` taxonomies, so a landing appears on the corresponding taxonomy
  term pages. Use the vocabulary already present on the site.
- `/projects/` lists landings alphabetically. The `projects/list.html` layout
  renders each card with its status and tags, linking them to their taxonomy
  term pages; the section index carries no introductory prose.
- Reusable `tools/rules_*` modules publish their documentation through
  `//tools:docs`, using external labels such as `@rules_example//:docs`. They
  are outside the landing registry and have no landing page, Pages repository,
  or DNS record. Keep standalone modules independent of the parent repository
  when integrating their documentation.
- Write the landing page for visitors, not as a copy of the README: a short
  description of what the project does, its status, and how to use it. Do not
  carry over repository-specific content such as source-tree links or
  repository build commands. Link to the project's documentation under
  `/docs/projects/<name>/` where it exists.
- When a project owns real user-facing documentation, it lives in
  `projects/<project>/site/content/docs/`, packaged with the same
  `docs_filegroup` recipe and depended on by the landing package. Its
  `_index.md` declares `cascade: [type: docs]` so the section and its subpages
  render with the documentation layout (section navigation and page layout)
  rather than the landing layout. Do not create a documentation directory for a
  project whose README is already its only user-facing documentation.
- Add the site's clickable link to the landing registry row in
  `projects/README.md`. Include a useful description and introduction in the
  project README. Use the shared Docsy header; do not duplicate navigation as
  buttons or add generated-source labels.

### Check the candidate

Use `repo-bazel` and `bazel-agent` for the owning build/tests and formatter.
Build `//projects/alwaldend.com:site`, check `//infra/dns:config_test`, and
inspect representative HTML from the candidate: the `/projects/` section
index, the landing page, its documentation page, and any taxonomy term page
it should appear on. Check its project identity, useful description,
documentation/source links, relative links and images, and preserved code
examples. Because the landing shares one build with the whole site, a
`--panicOnWarning` failure anywhere blocks the blog deployment. A successful
build alone does not establish a usable or deployed page.

### Deploy when authorized

Adding a landing page authorizes implementation and offline checks. Live
publication requires the user's requested operation and scope; retain
authorization already given instead of requesting it again. Publication is the
apex site deployment: `//projects/alwaldend.com:deploy` publishes the rendered
site to the apex Pages repository, following the
[Alwaldend.com README](../../../alwaldend.com/README.md) rollout entry points.
There is no separate landing deployment, and no per-project repository,
`bootstrap_projects` entry, Pages configuration, or DNS record to manage.

Plan DNS changes through the owning Terraform root, following `repo-infra` and
the [DNS workflow](../../../../infra/dns/README.md). Review the complete plan
and apply its saved artifact only within the authorized scope. Resolve any
unrelated infrastructure or apex changes before executing it.

Verify the pushed revision, Pages build state, and HTTPS page content. Report
observed successes separately from pending propagation or failed operations.
A Terraform apply or Git push alone is not evidence that the public site is
serving.
