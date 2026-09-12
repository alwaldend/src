## Context

See proposal.md — Why. This record evaluates how `alwaldend.com` should
distribute its landing page, generated documentation, and blog across
hostnames. The current state is one Hugo build published to the apex through
the `alwaldend/alwaldend.github.io` Pages repository.

Verified at source revision `387aed3b7e77c8bd850b8dc38ff88f92982ad50f`,
observed 2026-09-12:

- **The current structure is a deliberate consolidation.**
  `a535c248354f4ae0ca2a64aeb60b85cb92559321` (2025-06-23, "Rename hugo rules,
  combine misc and docs sites into one") merged the previously separate
  `hugo/sites/misc`, `hugo/sites/docs`, and `hugo/sites/projects` builds. The
  preceding `hugo/sites/misc/hugo.toml` carried a commented
  `baseURL = "https://alwaldend.com/misc"`, so even that earlier separation was
  path-based rather than hostname-based.
- **`/docs/` is a projection of the monorepo, not standalone content.** The
  site's `repo_docs` target aggregates `docs_filegroup` outputs from `//data`,
  `//infra`, `//projects`, `//third_party`, `//tools`, and `//users` under
  `prefix = "content/docs/"`, renaming `README.md` to `_index.md`. At this
  revision that tree packages 471 tracked `README.md` files (197 `projects`,
  164 `tools`, 51 `infra`, 41 `third_party`, 11 `data`, 6 `users`, 1 root).
- **The documentation tree is coupled to the blog build.** The site passes
  `--panicOnWarning` for both release and preview configurations, so a warning
  in any packaged README fails the same build that renders the blog.
- **The blog is currently one post,** `content/blog/slop-without-a-clear-goal`.
- **Hostname precedent exists and is cheap.** The apex zone already carries 29
  `CNAME` records pointing at `alwaldend.github.io.`
  (`infra/dns/zones/alwaldend.com!global.zone`). Pages repositories and their
  custom domains are provisioned declaratively by
  `infra/github/tf/alwaldend_pages_repos.tf`, which derives each landing
  hostname by replacing underscores with hyphens. Adding a hostname is
  therefore ordinary IaC work rather than a novel capability.
- **The reusable landing macro cannot host the docs or blog trees.**
  `al_hugo_landing`'s README rule strips shortcodes with
  `sed 's/{{[<%].*[>%]}}//'`, its generated config sets
  `disableKinds = ["taxonomy", "term"]`, and it emits only `home` and `page`
  HTML with no RSS or print output. The apex site's Docsy section outputs,
  taxonomies, link-resolution partials, and print-anchor scoping would have to
  be rebuilt for any extraction.
- **Cross-tree references are bounded.** 37 tracked references to
  `alwaldend.com/docs` exist across 23 files. Thirty of them sit in the 19
  `projects/*/landing/BUILD.bazel` files; the rest are the landing macro that
  defaults `docs_url` and a few tree READMEs.

## Goals / Non-Goals

**Goals**

- Record the evaluated hostname structures and the decisive evidence.
- Record the decision, the accepted trade-offs, and the reopening conditions.
- Give a future implementer the supported fallback for the one real cost.

**Non-Goals**

- Changing any hostname, DNS record, Pages repository, build target, or content.
- Restating `project-alwaldend-com` or `blog-section` requirements, which are
  unchanged.
- Establishing ranking, traffic, or search-visibility outcomes, which this
  record does not measure.

## Decisions

### Keep one hostname and one site build

`alwaldend.com` retains `/`, `/docs/`, and `/blog/` in the single existing
build and Pages repository. No `docs.` or `blog.` hostname is introduced.

The generated documentation tree is the deciding factor. It is a rendering of
the repository's own directory structure spanning 471 README files; giving it
its own hostname would not decouple it from that structure, it would only move
where the projection is served while requiring the Docsy section outputs,
resolution partials, and print scoping to be rebuilt outside the reused
landing macro. The prior consolidation points the same way for the same
sections.

### Treat the shared-build failure coupling as a build concern, not a hostname concern

`--panicOnWarning` currently makes any packaged README able to fail the blog's
build. That is a genuine cost, and it is the strongest argument for splitting.
It is addressed by separating build graphs rather than hostnames: two Hugo
builds over the same content — landing and docs at `/`, blog at `/blog/` —
merging their outputs into the one published tree. This yields independent
failure domains and release cadence while preserving one hostname, one
canonical identity, one client-side search index, and the existing link
structure. It is the recorded fallback if the coupling becomes painful.

### Do not separate the hostname for SEO reasons

Search-engine treatment does not decide this question at the margin, and this
record does not claim that a subdirectory preserves link equity that a subdomain
would forfeit.

Google Search Central's crawling-and-indexing FAQ answers the question
directly, and it is the strongest source here because it is edited vendor
documentation rather than transcribed speech. Under "Is it better to use
subfolders or subdomains?" it states: "You should choose whatever is easiest
for you to organize and manage. From an indexing and ranking perspective,
Google doesn't have a preference."
(`https://developers.google.com/search/help/crawling-index-faq`, retrieved
2026-09-12.) The vendor's stated criterion is therefore organizational and
operational convenience, which is exactly the consideration this decision
rests on.

A Google Search Central office-hours answer, checked 2026-09-12, gives the same
position in more detail. It is transcribed speech rather than edited vendor
documentation, so the quotation below may vary slightly between transcriptions.

- In the hangout published 2018-05-18
  (`https://www.youtube.com/watch?v=kQIyk-2-wRg`), Mueller said "we see these
  the same," that he "would personally try to keep things together as much as
  possible" and would "use subdomains where things are things are really kind of
  slightly different" (the duplication is in the published transcription), and
  that the choice "could go either way." (Search Engine Journal transcribes this
  answer at the 11:18 mark, embedding the same hangout:
  `https://www.searchenginejournal.com/google-treats-subdomains-subdirectories-john-mueller-says/254687/`,
  published 2018-05-25.)

None of these sources measures this site, and no traffic, ranking, or
search-visibility outcome is asserted here. The guidance also does not establish
that separating sections onto subdomains would be pointless — only that Google
does not require it and does not rank one form above the other. The
search-engine question is therefore settled as "acceptable either way," which
leaves build ownership and operational surface to decide it.

A third vendor video, "Subdomain or subfolder, which is better for SEO?"
(`https://www.youtube.com/watch?v=uJGDyAN9g-g`, published 2017-12-21), poses
the same question in its title. Its spoken content could not be retrieved when
this record was written — its captions were unavailable and the secondary
transcription sites were reachable only through bot-protection interstitials —
so no claim about what it advises is recorded here.

## Risks / Trade-offs

- **Accepted: the shared build remains a single failure domain.** A warning in
  any of the 471 packaged README files can fail the build that renders the
  blog. Accepted because the fallback above addresses it directly and no
  hostname change is required. Revisit if it blocks delivery in practice.
- **Accepted: sections share one visual and navigation surface.** A future
  blog needing a distinct theme or a CMS would have to override the Docsy shell
  or revisit this decision.
- **Accepted: one client-side search index.** `offlineSearch` indexes the whole
  site; a split would have produced two indexes with no cross-search. Keeping
  one site preserves unified search.
- **Risk: the documented hostname convention is not universal.** The project
  tree documents that `<hostname>.alwaldend.com` denotes a project landing page
  with hyphens in hostnames and underscores in documentation paths. `docs` and
  `blog` are not projects, so introducing them would require an explicit
  carve-out. Avoided here by not introducing them.
- **Risk: a future extraction must rebuild the section machinery.** Any later
  move of docs or blog must reimplement resolution partials, print scoping,
  taxonomies, and section outputs rather than reusing `al_hugo_landing`. The
  recorded fallback avoids this because both builds share one Hugo site
  definition, so the section machinery is reused rather than reimplemented.

## Migration Plan

No migration. The current structure is retained unchanged, so there is nothing
to deploy, roll back, or verify beyond this record itself.

## Open Questions

None. The deferrable question — whether to separate build graphs — is recorded
as the fallback with its trigger condition rather than left unresolved.

## Reopening conditions

Reconsider this decision if any of the following becomes true:

- The shared build's failure coupling repeatedly blocks landing or blog
  delivery. Apply the recorded fallback first.
- The blog grows enough to need a CMS, a distinct theme, or a release cadence
  that must not depend on the monorepo build.
- Documentation acquires a genuinely separate audience or contributor set with
  its own publishing workflow.
- A measured search-visibility or traffic problem is attributed to site
  structure with evidence distinguishing it from content and linking.

Absent one of these, relitigating the split requires new evidence rather than
reconsideration of the same points.
