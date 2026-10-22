## 1. Define browser acceptance and shared inputs

- [ ] 1.1 Write browser scenarios before implementation for navigation, history, embedded instances, unsafe filenames, empty/missing listings, unavailable servers, and theme modes; verify coverage against the browser specification.
- [ ] 1.2 Identify canonical project/version metadata and add the single configurable download endpoint; verify production uses `download.alwaldend.com` and fixture overrides do not alter production configuration.

## 2. Implement the reusable browser

- [ ] 2.1 Add the shared JavaScript module and Hugo partial/shortcode; verify two independently scoped instances render native Nginx JSON without duplicate component implementations.
- [ ] 2.2 Add `/downloads/` and the header entry using the existing site shell; verify light/dark rendering, responsive navigation, and keyboard access in the browser fixture.
- [ ] 2.3 Implement query-based directory navigation, direct file links, boundaries, and request states; verify refresh/history, encoded filenames, errors, and no credentialed requests.
- [ ] 2.4 Embed listings in the shared release partial using canonical metadata; verify each release remains scoped to its own directory and unavailable listings do not break release content.

## 3. Package and deploy the website through SSH

- [x] 3.1 Package the complete release-mode Hugo output with root `index.html`, required assets, and project landing pages; verify the generated archive contents.
- [x] 3.2 Wire the archive to the release tool's SSH site deployment and explicit environment selection; verify the configured targets and documented invocation use the generated bundle.
- [ ] 3.3 Update the production deployment entry point and documentation while preserving preview and unrelated staging; verify build/preview do not publish and staging inputs remain valid.
- [ ] 3.4 Extract the actual generated website archive into the fixture and verify representative page, asset, and project-landing requests.
- [ ] 3.5 Upload and select a concrete generated website release on the isolated host; verify the served bytes match the bundle.

Progress reconciled on 2026-10-22: commit
`04ddc88e90434568d66e480967587bc6630b1d86` records the stamped website-bundle
build and archive inspection. The release-only archive is declared in
`BUILD.bazel`, and the owner README documents `infra/download:publish.local`
and explicit Yandex selection. Tasks 3.4 and 3.5 preserve the uncompleted
acceptance portions of the original compound tasks; generic release-tool
fixtures do not prove acceptance of this concrete website bundle. The browser
and public cutover remain unfinished. No new build or publication ran during
this task-state reconciliation.

## 4. Verify and coordinate cutover

- [ ] 4.1 Build the site with the live download service unavailable and run relevant checks; verify the archive and listing assets are produced without infrastructure dependencies.
- [ ] 4.2 Run browser acceptance against the actual Nginx fixture; retain screenshots, request outcomes, candidate identity, and repeatable commands for Downloads and release-page instances.
- [ ] 4.3 Coordinate separately authorized site publication and apex DNS cutover with the hosting/DNS owners; verify both split-horizon destinations serve the same intended main hostname and preserve `www` TLS/redirects before retiring the production Pages dependency.
