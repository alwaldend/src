# Evidence manifest

[Back to goal](README.md)

## Protected baseline

| File | SHA-256 |
|---|---|
| `projects/renders/blender/fumo/reimu_fumo/reimu_fumo.blend` | `489213b7d0a62feb5c6b60ce36483757638886af3a4af25efa41e402e46b1d76` |
| `projects/renders/blender/fumo/fumo_sisyphus/fumo_sisyphus.blend` | `c5bd58ed9b29a6d67c398136eaec7ed34e227934c464662dfcb61f61f8e6f591` |

## Attempt 1

| Evidence | Value |
|---|---|
| Candidate | `out/fumo_concrete_drop_scaffold/attempt_01/fumo_concrete_drop_scaffold.blend` |
| Candidate SHA-256 | `86630e599525e40663ad01e4bd8f4c5c5f12e9cb127740440a7bbe501b77d292` |
| Contact frame | `22` |
| Minimum bottom Z | `-.00068653 m` |
| Late-motion span | `0 m` |
| Contact-sheet SHA-256 | `df4d3306329de260173b7271eaed50541773cf5a761a06ccc45d3b59058355a5` |
| Animatic SHA-256 | `5655461b83b78acbccb14c5c772759461f41441b412398fadd8138ef6a45a8cc` |
| Verdict | Technical pass; visual reject. |

The first MCP run stopped before saving on Blender 5.2's layered Action API.
The second reached rendered frames but found no Blender `FFMPEG` output enum.
The third completed after removing the nonessential legacy F-curve walk and
moving sampled-frame encoding to the host `ffmpeg` MPEG-4 encoder.  Neither
recovery changed the frozen physics or shot hypothesis.

## Attempt 2

| Evidence | Value |
|---|---|
| Candidate | `out/fumo_concrete_drop_scaffold/attempt_02/fumo_concrete_drop_scaffold.blend` |
| Candidate SHA-256 | `a9488c220c5076a3202e61c9897cf3710f24b1abe74fb9edfc4750bfaebfdc26` |
| Contact frame | `22` |
| Minimum bottom Z | `-.00068653 m` |
| Late-motion span | `0 m` |
| Contact-sheet SHA-256 | `23c5f9b5a879202ada16bfb3533b177d529e6289817c7b94cf384904e3d7f052` |
| Animatic SHA-256 | `490a4a25c47608e246bbb200a6ae1a94a9238a2702e3ef38b8131dd843ac2314` |
| Builder SHA-256 | `7558abd0406e085743616dc020fef44ceffc455d86981003633e74cb91c57799` |
| Build-report SHA-256 | `0971168f1cf21f80851896b142a13b51fb7954e168928b8761e27b5798288d8c` |
| Scripted checks | `11/11` pass |
| Verdict | Mechanics pass; composition reject; clean reopen unverified. |

The contact sheet was reviewed directly after hashing. It contains the neutral
warning and scale witness, but frames `1` and `12` crop most of the falling
proxy. The build report limits its claim to a rigid scene scaffold and records
both protected blend hashes unchanged.
