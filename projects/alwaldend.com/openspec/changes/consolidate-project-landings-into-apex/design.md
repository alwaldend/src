## Context

See proposal.md — Why. Verified in this checkout at the rebased base
`83dfca7f7b66c6b325e0f18a44d8c2b6e282ae2a`, observed 2026-09-19:

- `projects/projects.bzl` registers 19 landing projects. Each owns
  `projects/<project>/landing/BUILD.bazel` (19 files) using `al_hugo_landing_site`,
  which renders one `_index.md` from the project README through
  `sed 's/{{[<%].*[>%]}}//'`.
- The apex site already packages the same README at `/docs/projects/<name>/`:
  `//projects:docs` aggregates every project's `docs_filegroup`, and the apex
  renames packaged `README.md` files to `_index.md`
  (`projects/alwaldend.com/BUILD.bazel`).
- Publication surface per project: one Pages repository and one `CNAME`
  (20 catalog entries carry `landing_project`), one project `dnsconfig.json`
  (20 files; 19 landings plus the apex), and one Terraform DNS stage (20
  `projects/*/tf/dns.tf`). The global zone carries 29 `IN CNAME
alwaldend.github.io.` records.
- The apex packages an explicit source list (`site_source_filegroup`), so
  content inclusion is a declared input rather than a directory scan.
- The apex runs `--panicOnWarning` for release and preview, so any packaged
  content warning fails the build that also renders the blog. This change
  accepts that coupling deliberately as part of the monorepo build.
- Taxonomy terms are separate facts from taxonomy pages: a landing page
  declaring `statuses`, `languages`, or `tags` both renders term metadata and
  makes the apex emit the corresponding term page. The isolated landing builds
  set `disableKinds = ["taxonomy", "term"]`, so this participation is new.
- `al_hugo_landing` sets `disableKinds`, emits only `home` and `page` outputs,
  and produces exactly one page, so it cannot host documentation subpages.
- Aggregated documentation is non-recursive in most projects: `//projects:docs`
  deps on `//projects/<name>:docs`, and each project globs its own package level.
  Only `projects/tf_modules` and `projects/kustomization` recurse through
  `subpackages()`, so a new `site/` package is excluded by default there too and
  needs an explicit exclusion only if it declares a `docs` target.

## Goals / Non-Goals

**Goals**

- Give registered projects a visitor-facing landing page and optional
  documentation inside the main site, with project-owned content.
- Derive landing membership from one registry so a new project needs no second
  membership edit.
- Remove the per-project repository, custom domain, DNS record, and deployment
  path from the supported structure.

**Non-Goals**

- Rewriting project READMEs or moving repository documentation. READMEs keep
  their existing `/docs/projects/<name>/` URL and owner.
- Providing per-project visual theming. All landings use the shared shell.
- Deleting live DNS records, Pages repositories, or certificates. Retirement
  requires separately authorized infrastructure work.
- Deciding the fate of pre-existing drift discovered below.

## Decisions

### Separate the public namespace from the on-disk content directory

Publish at `/projects/<name>/` while packaging project content from
`projects/<name>/site/content/`. The public namespace is a Hugo section, which
is independent of the directory the build packages. Alternatives considered:
naming the content directory `projects/<name>/projects/` (duplicates the
segment and collides visually with the repository documentation projection at
`content/docs/projects/<name>/`), and renaming the public path to a short
segment such as `/p/` (saves characters, loses the self-describing name, and
departs from the documented project-naming vocabulary without adding a
section-level benefit).

### Use the registry as the single owner of landing membership

`projects/projects.bzl` gains the landing content package per project (or a
parallel declaration in the same file), and the apex site's packaged sources
derive from it. Alternatives considered: one hand-maintained list per consumer
(duplicates membership and desynchronizes), and a directory scan
(`glob`/`subpackages`) over `projects/*/site` (couples packaging to path
presence and cannot express a project with no landing yet).

### Classify pages by declared type, not directory position

A landing page declares the landing type and documentation pages declare the
documentation type, with the documentation index cascading the type to its
children, mirroring the existing `content/docs/misc` cascade. Alternatives
considered: relying on Hugo's implicit top-level-section type (every page under
`content/projects/<name>/` would resolve to the same type, so a landing and its
documentation could not render distinguishable layouts), and deriving layout
from a front matter `layout` value per page (works but duplicates the type
decision into every file and does not cascade).

### Keep landing content distinct from repository documentation

The landing is a short, hand-written visitor page; the README remains the
contributor-facing document. Alternatives considered: generating the landing
from the README as the retired macro did (two renderings of one source, and the
transformation cannot reliably rewrite repository-specific command examples),
and rendering one page that serves both audiences (conflates two audiences,
which is the problem this change addresses).

### Let landing pages participate in the site's taxonomies

Landing front matter keeps feeding `statuses`, `languages`, and `tags`, so the
site emits term pages and projects become discoverable by those terms. Accepted
deliberately: this is new behavior relative to the isolated builds and it adds
the landing pages to the site's taxonomy and sitemap surface.

### Guards for the newly coupled build

Because the landings join a build that runs `--panicOnWarning`, the project
content packages are validated together with the rest of the site, and the
apex build is the check that catches a broken link or shortcode before
publication. Accepted trade-off, recorded under Risks.

### Merge only what the main site still needs from the reusable landing project

`projects/hugo_landing` currently owns two different things: presentation the
main site already consumes, and machinery for standalone landing publication
that this change retires.

Keep, moving ownership into the apex site's own tree:

- The shared canvas and accent styles (`_shared_theme.scss`,
  `_shared_accent.scss`). The apex site already imports these through
  `//projects/hugo_landing/assets/scss:shared_theme`, so relocating them
  preserves one implementation instead of creating a second copy.
- The landing page rules (`.landing-content` sizing, image scaling, and
  table/pre overflow handling), which the apex needs for the new landing pages.

Do not port the reusable project's taxonomy and term layouts. They render a
bare heading, while the Docsy theme the apex already mounts provides full
taxonomy term clouds and term pages that list pages with descriptions,
breadcrumbs, and metadata. The apex has taxonomies enabled and no local
taxonomy or term layouts, so it would take the theme's richer defaults; porting
the minimal versions would be a behavioral downgrade.

Remove as having no remaining consumer: the landing macro
(`pkg/bzl/al_hugo_landing.bzl`) and its per-landing configuration generation,
the standalone landing configuration, the Go landing publisher (`cmd/deploy`),
and the repository-relative link and image rewrite layouts, which exist only to
point README links and images at GitHub source paths. In-site landing content
links within the site, where the apex's own resolution partial and bundled
content apply.

Alternatives considered: keeping `hugo_landing` as a reusable project with no
consumer (dead weight, and a second declared owner for shared styles), and
deleting it entirely (would delete the shared canvas and accent the apex
imports, forcing a copy).

## Risks / Trade-offs

- **Every landing content error can fail the whole-site build and the blog
  deploy.** → Accepted as the monorepo trade-off. Mitigation: keep landing
  content small and validate with the apex site build before publication.
- **A project's landing and its repository documentation can drift apart.** →
  Different owners and purposes are intentional; the landing spec requires the
  two sources to stay distinct rather than duplicated.
- **Retiring the subdomains breaks published links.** → The former hostnames
  stop resolving to a site. Mitigation: decide redirect handling during
  teardown; a per-hostname redirect page is cheaper than removing the
  repositories outright.
- **Repository deletion is blocked by design.** → `prevent_destroy = true` in
  `infra/github/tf/alwaldend_pages_repos.tf` protects adopted repositories, so
  teardown must change lifecycle and `landing_project` catalog membership
  deliberately rather than by removing a record.
- **Landing content does not exist yet for most projects.** → Landing pages are
  authored per project; a registered project without content produces no page
  and no failure. Documentation appears only where real user-facing
  documentation exists.
- **Pre-existing drift unrelated to this change.** → The historical global zone
  snapshot lists 11 `rules-*` hostnames whose landing infrastructure the module
  specifications already record as retired, and the catalog retains a
  `goal-landing` entry and repository for a removed `goal` project. Both predate
  this change; report them rather than silently folding them into retirement
  scope. The zone file is a snapshot and must not be hand edited.

### Scope retirement to each project's own DNS

Each landing project's Terraform root declares only DNS: 19 of the 20 project
`tf` directories contain `dns.tf`, `dns_variables.tf`, and `provider.tf` with no
other resource, so they are removed whole. `projects/alwaldend.com/tf` is the
exception — it also declares VM resources (`vms.tf`) — so it stays and keeps its
`pages` declaration, which the apex publication depends on. The project
`dnsconfig.json` declarations, and the `al.lua`, `al_config`,
`vault_binary_map`, and `dnsconfig` filegroup plumbing that served only those
stages, are removed with them.

Vault configuration is explicitly out of scope. The `src_projects_*` DNS
AppRoles, their modules under `infra/vault/tf/approles/`, and the
`group_dns_approles` membership are left untouched even though the stages that
consumed them are gone; retiring them is separate infrastructure work with its
own authorization.

The `tools/rules_*` modules need no removal work: their `project-dns`
specifications already record that landing infrastructure was retired, and they
carry no `dnsconfig.json`, Terraform root, or operational source exports. Their
specifications are verification targets, not work.

The `zones/` provider snapshots are removed rather than hand edited: they
duplicated the declarations and could not be kept current by inspection. They
are replaced by generated declaration pages, one per destination view, that
project the checked-in `dnsconfig.json` files and are verified for freshness by
`//infra/dns:config_test`. The pages are generated by `//infra/dns/cmd/dump`;
because a formatter could rewrite them into a state the generator would not
produce, they carry a `rules-lint-ignored` attribute alongside their freshness
check.

`projects/hugo_landing` also owned the per-project landing publisher and the
repository-relative rewrite layouts. These are removed with the standalone
configuration and the landing macro. Two leftover apex publication scripts,
`deploy_project.sh` and `deploy_all.sh`, published a landing repository
selected by project name and have no remaining caller; they are removed as
dead landing-publication plumbing. The `projects/BUILD.bazel` aggregate
landing targets and `//projects:deploy_landings` are removed with them.

The per-project DNS AppRoles under `infra/vault/tf/approles/`, including the
`src_projects_hugo_landing` module, are retired project infrastructure. They
are retained deliberately: Vault configuration is out of scope for this
change, so their retirement is separate work with its own authorization.

## Migration Plan

1. Land the in-site section and packaging first, with landings published at
   `/projects/<name>/` while the subdomains still resolve. Both structures are
   readable during the transition.
2. Author landing content per project, dropping repository-specific content
   such as source-tree links and repository build commands.
3. Verify the apex site build and its rendered `/projects/` output.
4. Retire dedicated publication: remove the landing targets and the
   landing deployment command and publisher, remove landing DNS declarations
   and stages, and remove landing repository and Pages membership from the
   repository catalog.
5. Rollback before step 4 is limited to reverting source changes; after step 4,
   restoring a subdomain requires re-creating its records and Pages
   configuration.

## Redirect handling

The user authorized the teardown and directed that the retired projects be torn
down rather than preserved. Removing the DNS records and the repositories is
what the authorization covers, and a redirect page would require retaining
both, so no redirect is published: a former hostname stops resolving. This
follows from the teardown authorization rather than a separate stated
preference, so record it as the chosen handling and revisit it before deleting
the records if a redirect becomes a requirement.

The live teardown remains a separately authorized infrastructure operation;
this change lands the source and catalog retirement that make it reconcilable.
It is additionally ordered after deployment: the apex site publishes through
the manual `//projects/alwaldend.com:deploy` step, which runs only after this
change merges, so `/projects/<name>/` is not serving while the subdomains still
are. Removing the DNS records and repositories before that deployment would
break every published project link with no replacement, so the live retirement
waits for the merged and deployed site.

The retirement source in this change is what makes the live teardown
reconcilable, and the two halves reach their targets differently. The catalog
removal makes `//infra/github/tf` plan the deletion of every landing repository
resource from its existing state, so that half is reachable from the delivered
tree once `prevent_destroy` is deliberately relaxed for the reviewed scope. The
per-project DNS records are not: their roots are removed by this change, so
their Vault-backed states have no owning target in the delivered tree. Deleting
those records means operating the pre-change project roots from revision
`83dfca7f` in a task-owned worktree, which is why the DNS half is an ordered
infrastructure step rather than part of the source change.
