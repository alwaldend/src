## 1. Landing registry and packaging

- [x] 1.1 Declare each project's landing content package in `projects/projects.bzl` and verify the registry loads with `bazel_agent bazel query //projects:all`
- [x] 1.2 Add the apex site's packaged source entry for project landing content, deriving membership from the registry rather than a second list, and verify `//projects/alwaldend.com:site_source_filegroup` builds
- [x] 1.3 Verify a project with no landing content produces no page and no build failure by building `//projects/alwaldend.com:site` while the registry includes such a project

## 2. Site section and classification

- [x] 2.1 Add the `/projects` section index page and confirm it renders at `/projects/` in the built output
- [x] 2.2 Add the landing page layout and front matter contract, and confirm a landing renders with the landing layout rather than the documentation layout
- [x] 2.3 Add the documentation type and section index cascade for `/projects/<name>/docs/`, and confirm a documentation subpage renders with the documentation layout under the project path
- [x] 2.4 Confirm landing pages render with the shared canvas styles and declare no project-local layout, style, or site configuration

## 3. Merge the reusable landing project

- [x] 3.1 Move the shared canvas and accent styles into the apex site's assets tree, repoint the apex imports at their new owner, and confirm the site still builds and renders the shared palette
- [x] 3.2 Move the landing page rules into the apex site's styles and confirm a landing page applies them
- [x] 3.3 Confirm the site renders taxonomy and term pages from the theme with no site-local taxonomy or term layout
- [x] 3.4 Remove the landing macro, the standalone landing configuration, the standalone landing publisher, and the repository-relative rewrite layouts, then confirm the repository builds without them

## 4. Landing content

- [x] 4.1 Author the landing page for each registered project as a short visitor-facing page, excluding repository-specific content such as source-tree links and repository build commands
- [x] 4.2 Add user-facing documentation under `/projects/<name>/docs/` only for projects that own real user-facing documentation, and confirm projects without it publish no documentation pages
- [x] 4.3 Confirm landing front matter feeds `statuses`, `languages`, and `tags`, and that the corresponding taxonomy term pages list the landing pages
- [x] 4.4 Confirm `README.md` still renders at `/docs/projects/<name>/` and that landing content is not published under the repository documentation path

## 5. Verification

- [x] 5.1 Build the apex site and inspect representative rendered HTML for the section index, a landing page, a documentation page, and a taxonomy term page, including relative links and images
- [x] 5.2 Confirm no landing page requires a per-project Pages repository, `CNAME`, or DNS record to be published in the main site

## 6. Documentation and skill alignment

- [x] 6.1 Update `projects/README.md` to describe the in-site `/projects/<name>/` structure and remove the per-hostname landing table
- [x] 6.2 Update the `repo-hugo` skill so its documented procedure matches the in-site structure instead of dedicated landing publication
- [x] 6.3 Update or remove the `project-hugo-landing` and `project-dns` specifications whose behavior is retired, keeping each fact owned once

## 7. Authorized retirement

- [x] 7.1 Remove dedicated landing deployment: the 19 project landing targets, `//projects:landing_sites`, `//projects:landings`, `//projects:deploy_landings`, and the landing deployment command, then confirm the repository builds without them
- [x] 7.2 Remove each landing project's DNS declaration, DNS stage, and the build plumbing that served only that stage (`al.lua`, `al_config`, `vault_binary_map`, and the `dnsconfig` filegroup), then confirm the repository builds without them
- [x] 7.3 Keep the apex `projects/alwaldend.com/tf` (it also declares VM resources) and its `pages` declaration; confirm the apex site still publishes and no Vault infrastructure is modified
- [x] 7.4 Confirm the `tools/rules_*` modules have no landing DNS declaration, Terraform root, or operational source export, and that their `project-dns` specifications still describe that retired state
- [x] 7.5 Remove landing repository and Pages membership from the repository catalog and handle `prevent_destroy` deliberately, then confirm the catalog precondition still passes
- [x] 7.6 Confirm `bazel_agent bazel test //infra/dns:config_test` passes with no landing hostname declared by any source, and that shared apex and mail declarations are unchanged
- [x] 7.7 Decide redirect handling for the former hostnames from the teardown authorization, and record the decision in the design
- [x] 7.8 Replace the provider-snapshot zone files with generated declaration pages per destination view, and confirm no landing hostname appears in a declaration
- [ ] 7.9 Tear down the former hostnames' live DNS records and the retired landing repositories, then verify the former hostnames serve no dedicated landing site. **Blocked, not performed**: the live check at 2026-09-19T12:52Z observed `https://alwaldend.com/projects/` returning 404 while every `<name>.alwaldend.com` still returned 200. The apex site publishes only through the manual `//projects/alwaldend.com:deploy` step, which runs after this change merges, so the replacement pages are not live yet. Tearing the subdomains down now would break every published project link with no replacement. Retire DNS and repositories only after the site deployment serves `/projects/<name>/` and `/docs/projects/<name>/`.

## 8. Validation and delivery

- [x] 8.1 Run the formatter and confirm the checked-in declaration pages are unchanged by it
- [x] 8.2 Confirm `//infra/src/openspec/validation:validate_test`, `//infra/dns:config_test`, `//:repo_quality_test`, `//.agents:write_skills_test`, and the affected site and skill targets pass
- [x] 8.3 Build the whole repository with `bazel_agent bazel build //...` and confirm it succeeds
