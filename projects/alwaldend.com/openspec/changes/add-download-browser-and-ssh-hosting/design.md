## Context

See [proposal.md](proposal.md). One Hugo/Docsy site owns the shared theme,
navigation, and release partial. The body-end hook currently packages
`assets/js/navigation.js`; release pages share
`layouts/_partials/alwaldend/release.html`. The current production deploy
script publishes the generated site to a GitHub Pages branch.

This change is the only active replacement of `Separate preview and
publication`. The sibling
[landing consolidation plan](../consolidate-project-landings-into-apex/proposal.md)
retains its build delta and the shared-publication requirement in its added
`project-landing-pages` capability, but no longer replaces the publication
requirement with the former Pages workflow. Either archive order therefore
preserves SSH publication and shared landing publication. Its completed-task
and deployment evidence remains historical evidence, not a new rollout step.

## Goals / Non-Goals

**Goals:** Render live file listings in the existing website using a small
reusable browser component, and publish the built site through SSH releases.

**Non-Goals:** Another Hugo site, a copied theme, a frontend framework, a
server-side browsing application, upload UI, or changes to unrelated staging.

## Decisions

### Shared listing component

Add one module under `assets/js/` and one Hugo partial, exposed through a
shortcode where content needs it. Each instance supplies an endpoint, starting
directory, and navigation boundary. Load the module once on pages using it
and permit multiple independent instances on a release index.

The `/downloads/` page starts at `projects/` and uses the existing site shell.
Add Downloads to the standard header menu. A `path` query parameter records
the selected directory so refresh, links, and browser history work with the
same static page. Directory navigation stays in the component; file links
point to the download host's actual file URLs.

Release pages pass their canonical project key and concrete release version
to the same partial. An embedded listing cannot navigate above its configured
release root. Keep project/version mapping owned by release metadata, not a
second frontend registry. Release indices can defer fetching collapsed
sections until opened to avoid loading every directory at once.

Nginx [native JSON listings](https://nginx.org/en/docs/http/ngx_http_autoindex_module.html)
provide names, entry types, dates, and file sizes. Render filenames as text,
encode URL path segments, validate navigation boundaries, and use the site's
existing typography and table styles. Show loading, empty, missing-release,
and unavailable-server states locally without breaking surrounding content.
Provide direct download-root links as a fallback when JavaScript is disabled.

Use credential-free cross-origin fetches from
`https://download.alwaldend.com`; the infrastructure owner configures CORS.
Keep this URL in one site configuration parameter shared by all instances.
For isolated previews, allow an explicit fixture endpoint override. Production
split-horizon DNS selects the environment without different frontend builds.

### Ownership and build boundaries

The UI and theme remain entirely in `projects/alwaldend.com`; Nginx only
provides data and file bytes. The reusable Nginx role has no UI source. The
site build must succeed with the download service unavailable and must not
depend on `infra/` production artifacts or live JSON responses.

### Production archive and deployment

Package the existing release-mode Hugo output as an archive with `index.html`
at its root. Include generated CSS, JavaScript, icons, and other referenced
assets and packaged project landing pages. Add it to the existing release
tool's manifest with an explicit SSH
site deployment and public project key `alwaldend.com`.

`//projects/alwaldend.com:deploy` becomes the explicit production SSH deployment
entry point after implementation; building and previewing remain read-only
with respect to hosts. Select local or Yandex through explicit deployment
inputs, not mutable DNS-based SSH targeting. Keep the unrelated staging
deployment working unless a later request changes it.

The public URL remains `https://alwaldend.com`. The hosting owner serves the
selected extracted release. Migrating apex A/AAAA records belongs to
`infra/dns`, while download records belong to `infra/download`. Do not put
duplicate apex declarations in the website or download component.
The hosting owner also preserves the existing `www` alias with trusted HTTPS
and a permanent redirect to the apex, retaining paths and queries.

## Risks / Trade-offs

- Download server unavailable -> show an explicit retryable listing error;
  the main page and already available direct links remain usable.
- Multiple listing instances -> isolate state and event handlers per instance;
  test two releases on one page with independent navigation boundaries.
- Arbitrary filenames -> render text rather than HTML and encode path segments;
  test characters that would otherwise become markup or query delimiters.
- JavaScript-free clients -> direct endpoint access remains available, but the
  themed directory browser requires JavaScript.
- New host plus new UI in one release -> verify the built archive against a
  fixture server before any production DNS transition.

## Migration Plan

Build and inspect the candidate archive and browser behavior using controlled
JSON fixtures and an isolated Nginx instance. Publish to explicitly authorized
VMs and verify both content hostnames and the `www` HTTPS redirect by address
override. Coordinate apex cutover
with its DNS owner only after those checks pass. Keep GitHub Pages reachable
through preparation; remove it as the production deployment dependency after
cutover. Changing the selected site release uses ordinary redeployment; no
rollback command is introduced.
