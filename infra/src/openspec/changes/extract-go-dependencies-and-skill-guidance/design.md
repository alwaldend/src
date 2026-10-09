## Scope and ownership

Copy only shared Go dependency configuration and agent skill changes from
33cc18d264fc36a5b0f9b075c7afa49352d8f524 onto freshly fetched master
4c03ac887c8a70e1533e765f4ff2245dcde2f5dc. Preserve the existing skill change
records with their owners. This record owns the extraction and dependency
delivery; the Anki project record continues to own CLI implementation.

## Dependencies and delivery

The approved SQLite v1.60.1 pin and its approved transitive dependencies remain
unchanged. The Go SDK version is unchanged. Additional module version changes
are the existing resolved dependency graph, not a fresh upgrade request.

Publish the independent prerequisite against master first. The project uses
SQLite's collation registration API, unavailable in master's v1.20.3 driver.
Removing its dependency changes before this prerequisite merges would prevent
validation. Preserve PR #128's published head until the prerequisite is merged,
then synchronize and validate its reduced diff before updating it. Do not merge
the prerequisite or retarget PR #128 without user authority.

## Evidence

Extraction equality, exact-candidate validation, and publication are recorded
under ignored out/anki-split. Delivery receipts own observed publication state.
