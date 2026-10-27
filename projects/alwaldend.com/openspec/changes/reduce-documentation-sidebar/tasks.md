## 1. Configure and verify

- [x] 1.1 Confirm pinned Docsy behavior and configure scoped owner roots.
- [x] 1.2 Build the release archive, measure output, and inspect representative navigation HTML.
- [x] 1.3 Verify browser navigation and existing site checks; document and publish the result.

## 2. Investigate first-interaction latency

- [x] 2.1 Inspect sidebar and homepage hover behavior and compare first/repeated interactions after reload, with and without site scripts. Record measured results and reproduction limits without claiming the reported delay is fixed.

## 3. Upgrade Docsy and verify compatibility

- [x] 3.1 Pin Docsy 0.18.0 and update canonical theme references and required compatibility settings.
- [x] 3.2 Build the site, run browser checks, and measure first interaction in headed Firefox private browsing.
- [x] 3.3 Record upgrade results and prepare the existing pull request for repository delivery.

## 4. Deliver website fixes independently and remove interaction blocking

- [x] 4.1 Isolate website changes, documentation packaging, and font subsetting from infrastructure changes on a branch based on master.
- [x] 4.2 Precompute the full Lunr index during the site build and load/query it in a worker, preserving fields, ranking, exclusions, and result metadata.
- [x] 4.3 Verify search results, first interactions, responsive navigation, tags, archive contents, and deployment-independent builds.
- [x] 4.4 Document current behavior and deliver the website-only pull request.
