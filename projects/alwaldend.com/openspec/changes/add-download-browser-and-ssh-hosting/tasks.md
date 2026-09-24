## 1. Define browser acceptance and shared inputs

- [ ] 1.1 Write browser scenarios before implementation for navigation, history, embedded instances, unsafe filenames, empty/missing listings, unavailable servers, and theme modes; verify coverage against the browser specification.
- [ ] 1.2 Identify canonical project/version metadata and add the single configurable download endpoint; verify production uses `download.alwaldend.com` and fixture overrides do not alter production configuration.

## 2. Implement the reusable browser

- [ ] 2.1 Add the shared JavaScript module and Hugo partial/shortcode; verify two independently scoped instances render native Nginx JSON without duplicate component implementations.
- [ ] 2.2 Add `/downloads/` and the header entry using the existing site shell; verify light/dark rendering, responsive navigation, and keyboard access in the browser fixture.
- [ ] 2.3 Implement query-based directory navigation, direct file links, boundaries, and request states; verify refresh/history, encoded filenames, errors, and no credentialed requests.
- [ ] 2.4 Embed listings in the shared release partial using canonical metadata; verify each release remains scoped to its own directory and unavailable listings do not break release content.

## 3. Package and deploy the website through SSH

- [ ] 3.1 Package the complete release-mode Hugo output with root `index.html`, required assets, and project landing pages; verify archive contents and representative page/asset/landing requests after fixture extraction.
- [ ] 3.2 Wire the archive to the release tool's SSH site deployment and explicit environment selection; verify a concrete generated release can be uploaded and selected on the isolated host.
- [ ] 3.3 Update the production deployment entry point and documentation while preserving preview and unrelated staging; verify build/preview do not publish and staging inputs remain valid.

## 4. Verify and coordinate cutover

- [ ] 4.1 Build the site with the live download service unavailable and run relevant checks; verify the archive and listing assets are produced without infrastructure dependencies.
- [ ] 4.2 Run browser acceptance against the actual Nginx fixture; retain screenshots, request outcomes, candidate identity, and repeatable commands for Downloads and release-page instances.
- [ ] 4.3 Coordinate separately authorized site publication and apex DNS cutover with the hosting/DNS owners; verify both split-horizon destinations serve the same intended main hostname and preserve `www` TLS/redirects before retiring the production Pages dependency.
