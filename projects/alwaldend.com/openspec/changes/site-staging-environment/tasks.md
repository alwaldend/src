## 1. Declarations

- [x] 1.1 Add the `www-staging` DNS declaration to `infra/dns/dnsconfig.json`
- [x] 1.2 Add the `www-staging` catalog repository and the creation-only `auto_init` default
- [x] 1.3 Add the staging site target and staging deployment arguments
- [x] 1.4 Regenerate the DNS declaration pages and confirm `//infra/dns:config_test` passes
- [x] 1.5 Confirm the catalog and Terraform formatting tests pass

## 2. Landing index

- [x] 2.1 Remove the explanatory header paragraph from the `/projects/` section index
- [x] 2.2 List landings alphabetically by title
- [x] 2.3 Render each card's status and tags, linked to their taxonomy term pages
- [x] 2.4 Document the index's ordering and card contents in the `repo-hugo` skill
- [x] 2.5 Give the section index its own `projects_index` type so it renders the
      centered documentation column without either sidebar, matching `/docs/`
      and `/blog/` geometry

## 3. Staging publication

- [x] 3.1 Apply the reviewed repository and DNS changes for the staging hostname
- [x] 3.2 Publish the site to the staging repository and confirm the Pages build succeeds
- [x] 3.3 Verify `/projects/`, `/projects/<name>/`, `/docs/projects/<name>/`, and the taxonomy pages on `https://www-staging.alwaldend.com/`

## 4. Production publication

- [x] 4.1 Deploy the accepted output to the apex site and confirm `/projects/<name>/` is served
- [x] 4.2 Verify the apex documentation and taxonomy routes still resolve
- [x] 4.3 Deploy the landing index changes to staging and the apex

## 5. Retirement (authorized 2026-09-20)

- [x] 5.1 Retire the per-project DNS records and confirm the former hostnames no longer serve a landing site. All 19 retired records destroyed; none of the former hostnames resolve.
- [x] 5.2 Retire the per-project landing Pages repositories and their GitLab mirrors. 20 GitHub repositories and 20 GitLab projects destroyed; the 9 shared repositories on each forge are unchanged.
- [x] 5.3 Record the retirement evidence and close
      `consolidate-project-landings-into-apex` task 7.9. Evidence is recorded there, including the plan-only scope, the branch-protection handling, and the removal of the unowned `goal.alwaldend.com` residual.

## 6. Validation

- [x] 6.1 Run the formatter and confirm no unrelated target changes
- [x] 6.2 Build the repository and confirm the affected site, DNS, catalog, and skill targets pass

## 7. Deferred

- [ ] 7.1 Sort the project index automatically by lifecycle rather than by
      title. Deferred at the user's request; the index currently orders by
      title.
