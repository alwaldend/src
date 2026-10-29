## Theme parity

The website's two generated CSS files match the prior candidate byte for byte.
All fourteen spatial model payloads, including resolved appearance data, are
unchanged. Stable browser screenshots match; native touch-scroll images retain
their timing-dependent offsets. Native SVG/raster checks pass with the JSON
palette independently supplying expected pixel colors.

The shared preset lives in the data tree because it is an ordinary diagram
input. An initial tools-tree placement failed transitive visibility; moving
the preset assets to their identity owner corrected the dependency boundary
without a new production exception. Generic provider/rule code remains under
`tools/mermaid`.

Final delivery includes the repository quality gate, semantic lint of the
aggregate affected packages, existing native-renderer tests, website browser
tests, spatial browser and snapshot tests, all sample snapshots, architecture
outputs, and the site build. Delivery receipts own exact candidate identities
and publication state. No new dependency or deployment was introduced.
