## Why

Every one of the 19 projects registered in `projects/projects.bzl` publishes a
dedicated landing site through its own GitHub Pages repository, custom-domain
`CNAME`, and per-project Terraform DNS stage, yet each renders a single page
derived from the project's `README.md`. The main site already renders that same
README at `/docs/projects/<name>/`, so the existing structure pays a
repository, a certificate, a DNS record, and a bespoke publisher for each
project to produce the weaker of two renderings of one file: the landing
strips Hugo shortcodes and emits no taxonomies, section outputs, or RSS.

The main site also conflates two audiences. `/docs/` projects the repository's
own filesystem layout, including contributor-facing READMEs and per-package
reference material, while offering no visitor-facing page for a project. This
change gives projects a visitor-facing home inside the main site and retires
the dedicated landing publications.

## What Changes

- Add a `/projects` section to the main Hugo build. Each registered project
  gets `/projects/<name>/` as a landing page and, when the project has
  user-facing documentation, `/projects/<name>/docs/`. The section index at
  `/projects/` lists the registered projects.
- Each project owns its landing content in `projects/<name>/site/content/`.
  The main site packages that directory at `content/projects/<name>/`. These
  directories carry pure content and no layouts, styles, or build rules; the
  theme and presentation shell remain shared.
- Landing content is a short, hand-written page for visitors, not a copy of the
  project README. The README stays owned by repository documentation and keeps
  its existing `/docs/projects/<name>/` URL, so contributor-facing reference
  material and visitor-facing pages have distinct owners and distinct sources.
- Landing front matter participates in the site's existing `statuses`,
  `languages`, and `tags` taxonomies. This is new behavior: the isolated
  subdomain builds had taxonomies disabled.
- Retire dedicated landing publication: the per-project landing targets, the
  landing deployment command, the 19 landing Pages repositories, their `CNAME`
  records, and their Terraform DNS stages. **BREAKING**: the existing
  `<name>.alwaldend.com` URLs stop being published.
- Keep `projects/projects.bzl` as the single owner of landing membership; the
  main site's packaged sources derive from that registry rather than a
  duplicated list.
- Replace the provider-snapshot zone files with generated declaration pages so
  the documented DNS inventory projects the declarations that own the records.
- Remove per-project DNS ownership: each landing project's `dnsconfig.json`
  declaration and its Terraform DNS stage, removing the whole Terraform
  directory where DNS was its only content. The apex Terraform root is retained
  because it also declares VM resources, and no Vault configuration is changed.

## Capabilities

### New Capabilities

- `project-landing-pages`: the `/projects` section of the main site, its
  registry-driven content packaging, the per-project content package
  convention, and the retirement of dedicated per-project landing publication.

### Modified Capabilities

- `project-alwaldend-com`: the packaged site build gains the project landing
  section as a packaged content input, and the site's published URL structure
  and taxonomy participation extend to it.

## Impact

- `projects/alwaldend.com/`: packaged site sources, site configuration, the
  `/projects` section index, and section layouts.
- `projects/<project>/site/content/`: new project-owned content directories,
  created only for projects with landing content to publish.
- `projects/hugo_landing/`: removed entirely. Its shared canvas and accent
  styles and its landing page rules move into the main site's tree; its landing
  macro, generated configuration, standalone publisher, and repository-relative
  rewrite layouts have no remaining consumer.
- `projects/<project>/landing/`, `projects/<project>/dnsconfig.json`,
  `projects/<project>/tf/`, and the landing catalog entries under
  `infra/repos/`: retired publication and DNS ownership.
- `projects/README.md` and the `repo-hugo` skill: document the in-site
  structure instead of the subdomain structure.
- `infra/github/tf` and `infra/repos`: the per-project Pages repository,
  `github-pages` environment, and landing default-branch resources are removed,
  and the landing repository records leave the catalog.
- `infra/dns`: the `zones/` provider snapshots are replaced by generated
  declaration pages, one per destination view, verified for freshness by
  `//infra/dns:config_test`.
- Live infrastructure: the landing Pages repositories, DNS records, and
  certificates require an explicitly authorized teardown.
- Merge the reusable `hugo_landing` project into the main site: keep its shared
  canvas and accent styles and its landing page rules, and remove the rest. The
  main site already imports those style files, so they change owner rather than
  multiply. Remove the landing macro, the generated landing configuration, the
  standalone landing publisher, and the repository-link rewrite layouts. Keep
  the theme's taxonomy and term layouts rather than the reusable project's
  heading-only versions, which would downgrade the term pages.
