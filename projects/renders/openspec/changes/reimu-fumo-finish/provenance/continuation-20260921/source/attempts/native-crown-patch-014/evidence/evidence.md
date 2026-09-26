# Native crown-patch reset receipt

## Exact model edit

The writer opened a task-owned copy of protected A157 whose SHA-256 remained
`433d08ad36be488bb16e4221a85f831d4390660c258a43ea0b08775811574b73`.
On `A44 continuous hair cap with smooth opening`, code established Edit Mode
face selection from the frozen `Z = 0.155 m` cut and Blender's native delete
operator performed the geometry change. No vertex coordinate was assigned.
The obsolete decorative `Subtle crown center seam` alone was hidden.

The saved candidate is
`out/reimu_fumo_finish/attempt_014_native_crown_patch/a157_native_crown_patch_014.blend`,
SHA-256
`fd6878d6e2dae076e005d937741df81e094d60087e284b7441b347db25a01724`.
The rejected bytes remain ignored task scratch.

Two pre-save failures were caught by invariants and produced no candidate.
First, opening a blend and entering Edit Mode in the same MCP request left an
invalid operator context; a later request separated file opening from editing.
Second, stale vertex and edge selection synchronized to every face on Edit
Mode entry; the corrected request explicitly deselected all BMesh elements
before selecting the 3,583 target faces. These were causal context and
selection repairs, not equivalent retries. The protected file was reopened
from unchanged bytes between them.

## Pinned audit

Repository-pinned Blender 5.2.1 clean-opened both files. The audit, SHA-256
`6cd9888f0545a3158117d10d4c553ebf46d03e8300ba42aea99939fe335cfb5e`,
found:

- changed geometry exactly `A44 continuous hair cap with smooth opening`;
- changed visibility exactly `Subtle crown center seam`;
- no added or removed objects, images, or linked libraries;
- unchanged modifiers and fixed cameras; and
- one remaining cap component with 5,252 faces, 5,350 referenced vertices,
  5,909 total vertices, and 194 boundary edges.

## Frozen fast trio

Pinned Blender 5.2.1 rendered the exact candidate at 512 by 512. Candidate
hashes before and after rendering matched. The fresh packet contains:

- front: SHA-256
  `054ff5b6a0a85c65cf7fc0aa10a6711c4737a7f767d4aa9b2b44968fcd4f04d3`;
- side: SHA-256
  `766895ddd881791593ac94843a24d04f7c84a4283059890115725610e255896d`;
- three-quarter: SHA-256
  `47a70bfc461398a60a54c4c08e165f50d1767a8446132dcc4d2586f7e7080f89`;
- manifest: SHA-256
  `86270e15a46eef5bb5cd77e95faee200518ec58d2ec2ea1c490641c7a1621488`.

Both coordinator and implementation-blind review returned `RESET`. Side and
three-quarter views show a broad pale band under a smooth brown cap terminated
by a straight horizontal rail. The remaining cap is a smaller helmet slab;
the fringe, side locks, cheek sawteeth, and lone rear lock read as rigid or
pasted-on pieces. The reference instead shows overlapping hair coverage and
softly rooted panels. These automatic construction failures make remaining
views and landmark measurement uninformative.
