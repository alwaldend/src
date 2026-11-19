---
title: Alwaldend blog skill evaluations
---

The offline target validates the Promptfoo configuration, case files, and skill
staging without credentials or model calls. Cases cover verbatim preservation
of supplied content, the unpublished default state, and reporting style
findings instead of silently editing them.

No live target is declared: representative behavior requires creating a Bazel
package, building the site, editing repository files, and reading rendered
output. Those effects cannot be exercised safely in this read-only staged
fixture. Configuration validation does not establish correct agent behavior or
a published post; live evidence belongs to the actual authorized task.
