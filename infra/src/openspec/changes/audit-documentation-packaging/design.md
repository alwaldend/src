## Context

The root docs aggregate reaches component documentation and generated release
pages. Those release pages depend on artifact metadata generation. The requested
website must contain documentation and site assets, with release browsing
provided by the download service instead of embedded release records.

## Decisions

Repair missing declarations at their owner. Remove release documentation edges
from project documentation aggregates, preserving release publishing targets.
Audit actual packaging output, not extensions in dependency names: build tools
and artifact metadata inputs do not prove that a binary is shipped in the site.

## Acceptance

Enumerate all root-workspace docs targets, build the relevant documentation
graph, inspect packaged paths, and inspect the final website archive. Fail on
generated release records or project download payloads; preserve site assets.
Record unrelated target failures and resolve those blocking the site build.

## Observed validation

All 556 docs targets pass analysis. The site build exposed two missing edges:
the gRPC dependency lacked docs, and the site cmd aggregate referenced a child
without docs. Repair both at their owners. Exclude 14 project release-doc
edges while keeping actual release targets. The completed archive contains no
project executable payloads or generated release records.

## Session findings

- `docs-payload-assumption`: dependency names initially suggested executable
  leakage; inspecting actual archive contents corrected that claim. Inventory
  the packaged output before inferring its payload from dependency names.
- `preview-snapshot-timeout`: the collaborative browser failed twice at the
  same snapshot step. The existing pinned browser workflow verified scoped
  navigation and produced a mobile screenshot and machine-readable receipt.
- `ssh-identity-environment`: OpenSSH -i checked an injected reference as a
  literal path. IdentityFile expansion and a real SSH fixture resolved it;
  authenticated local duplicate publication then refused as intended.
- `gazelle-unrelated-import`: full generation encountered an existing unresolved
  python.runfiles.runfiles import in repo_quality. Scoped Mermaid generation
  and the owning Python manifest generator completed; no unrelated owner edits.
