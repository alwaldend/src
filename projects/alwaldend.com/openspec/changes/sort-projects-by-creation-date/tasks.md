## 1. Project ordering

- [x] 1.0 Confirm historical creation dates against GitLab backup repositories
      using the existing infrastructure token and retain non-secret date evidence.
- [x] 1.1 Set explicit creation dates on all registered landing pages and use
      descending date order on the homepage and project index; verify both rendered
      lists against the front matter and confirm unchanged project membership.
- [x] 1.2 Require the creation date in the Hugo skill, update its existing
      onboarding evaluation, and document the ordering in the site README; verify
      the skill validation and offline evaluation configuration.

## 2. Validation

- [x] 2.1 Validate the OpenSpec record, build the site, run the relevant tests,
      mandatory repository quality and semantic lint checks, and retain rendered
      output inspection evidence.

Git and the AL preparation/publication receipt own delivery state. The receipt
is retained under `out/sort-projects-by-creation-date/` rather than duplicating
publication state in this implementation checklist.

## Current evidence

The site build, 35 selected tests (site, DNS, skill configuration/discovery,
Go login, and repository quality), semantic lint, and strict change validation
passed. The login command fixture covers 20 scenarios.
The built homepage and project index both contain all 24 registered projects
in descending declared date order. Representative documentation and taxonomy
links remain present. Local command logs and the ordering report are under
`out/sort-projects-by-creation-date/`.

Historical confirmation used full, non-shallow Git histories from the authorized
GitLab backups after the operator successfully refreshed Vault authentication.
Seven provisional dates were corrected. Android launcher, Autoscroll, CI
platform, Dotfiles, Nexus security plugin, and Useless QT GUI use their original
repository root commits. Leetcode downloader uses the first preserved source
commit from the older aggregate repository. Renamed projects were matched by
their retained implementation. Repository/API creation timestamps were not used.

Dotfiles contains a descendant with an author timestamp older than its root;
the root commit defines its creation date rather than that nonmonotonic timestamp.
Only dates, hashes, and project identity evidence are retained in task scratch;
credentials, authors, emails, and private source are not publication evidence.

Delivery uses `gh` as explicitly requested by the user. The merged AL transport
support is selected with `--git-transport gh`. GitHub reports write access to
`alwaldend/src`; the preparation receipt records current refs and PR identity.
