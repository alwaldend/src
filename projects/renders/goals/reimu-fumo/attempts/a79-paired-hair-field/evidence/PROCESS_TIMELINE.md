# A77--A79 process timeline and optimization

Snapshot: 2026-08-31 14:46 +03. Times below are local (+03). Canonical goal
timestamps, filesystem mtimes, exact artifact hashes, and render manifests were
used. Mtimes identify artifact completion, not necessarily worker start.

## Direct answer

The problem was pre-render process churn, not Blender. A77 and A78 consumed
58m 58s of canonical attempt time without candidate pixels. A79 then took
94m 30s from canonical start to its first whole-character pixel. Once invoked,
all five 512 px A79 views rendered in 29.115s and immediately exposed a bald
crown, rectangular curtain rear, planar board-like leaf, and broken joins.

The resulting sheet is useful evidence, but it is explicitly
`NON_CANDIDATE`, `promotion_authorized: false`, and a render of a rejected
frozen representation:

- `diagnostic_preview/a79_non_candidate_five_view.png`
- `diagnostic_preview/DIAGNOSIS.md`
- `diagnostic_preview/manifest.json`

The user's observation was therefore correct: until 14:41, there had been no
new whole-character A77--A79 geometry render to see.

## Elapsed timeline

| Time | Elapsed | Event and outcome |
| --- | ---: | --- |
| 12:01:11 | A77 +0 | Canonical A77 starts. |
| 12:15:17 | +14m 06s | A77 reference/baseline board appears. It is not candidate pixels: all five “A77 P0” cells point to the rung-003 baseline. |
| 12:31:41 | +30m 30s | A77 candidate `.blend` saved. |
| 12:38:26 | +37m 15s | Clean reopen/integrity passes, but mechanical preflight and render authorization fail; independent audit issues NO-GO. |
| 12:44:48 | +43m 37.5s | A77 closes `reset`, with no render. Face fit is 0.422 mm RMS / 1.237 mm maximum versus 0.25 / 0.45 mm limits. |
| 12:48:15 | A78 +0 | Canonical A78 starts after a 3m 26s gap. |
| 13:00:10--13:00:56 | +11m 55s--12m 42s | History/adversarial reviews reject the receiver-only boundary and cloth branch. |
| 13:03:35 | +15m 20.6s | A78 closes `reset` before geometry, Blender, or rendering. |
| 13:05:07--13:09:37 | post-close | Eight substantive A78 files still arrive over the next 6m 01s; only the interface inventory is reused by A79. |
| 13:06:37 | A79 +0 | Canonical A79 starts after a 3m 01s gap. |
| 13:10:34--13:14:26 | +3m 57s--7m 50s | 125 reference frames and reference contact/crop images appear. These are source pixels, not model progress. |
| 13:24:54 | +18m 18s | Variant B rejected pre-Blender as faceted/card-prone. |
| 13:35--13:49 | +28m--43m | Render harness, density analysis, registration, visual-hull analysis, and interface contracts complete in parallel. No candidate exists. |
| 13:55:30 | +48m 54s | `PROCESS_CHECKPOINT.md` explicitly says to stop expanding preflight and make the next artifact candidate pixels. |
| 13:57:26--14:04:37 | +50m 49s--58m 01s | Deferred standalone material branch produces two actual swatch renders, two enlarged derivatives, two swatch blend files, and a report. It does not touch the model. |
| 14:11:42 | +65m 05s | Separate rear-panel proposal rejects itself before Blender. |
| 14:14:44 | +68m 07s | Canonical interim checkpoint: A79 remains open; no candidate/render is authorized. Canonical state is now behind later scratch work. |
| 14:25:19 | +78m 43s | First frozen geometry receives a pure-preflight `PASS`; an exact independent replay finishes at 14:28:21. |
| 14:35:05 | +88m 28s | Variant A is reset after later checks find omitted bridge extents, rewritten inner controls, and interface incompatibility. |
| 14:39:07--14:39:23 | +92m 30s--92m 47s | Variant C resets; candidate scaffold audit remains `REVISE`/no authorized build. |
| 14:41:07 | +94m 30.0s | First new A79 whole-character pixel: front view of a deliberately non-candidate diagnostic. |
| 14:41:30 | +94m 54s | Fifth view completes. Manifest reports 29.115s total render time. |
| 14:43:05 | +96m 28.5s | Five-view sheet completes. From A77 start, the first new geometry sheet took 2h 41m 54s. |
| 14:45:16 | +98m 40s | Pixel diagnosis records a categorical representation failure. |

## Where the time went before pixels

By the first A79 whole-character render, the A79 tree held about 399 files and
148 MiB, including 2,909 Markdown lines and about 13,000 Python lines. Of the
218 PNGs, 125 were extracted source frames; most others were masks, overlays,
contact sheets, a rear projection, or material swatches. There was still no
authorized A79 candidate `.blend`.

Before the 13:55 process checkpoint, A79 produced 317 files, eight Markdown
documents, 13 Python files, and 207 PNGs. Between that explicit “pixels next”
decision and the first geometry pixel, it produced another 64 files, including
13 Markdown documents and 11 Python files. The critical-path instruction was
recorded but not enforced.

Rendering itself was cheap: the diagnostic script's first artifact to first
pixel took about 80.4s, and the five render calls totaled 29.115s. A pure-Python
frozen geometry already existed at 14:25:19; having the disposable preview
consumer ready then would have returned the same visible veto roughly 14
minutes earlier. Starting a crude preview at the 13:55 stop decision could
have avoided up to about 44 minutes, although that larger saving is an
estimate because the final frozen geometry did not yet exist.

## False-pass causes

No visual candidate was falsely accepted. The failure was that narrow passes
were easy to read or consume as global progress:

1. A77 reports `clean_reopen_pass: true` while also reporting mechanical
   preflight false and render authorization false. Its “A77 P0” board contains
   only byte-identical rung-003 baseline images.
2. A78 reports topology preflight `PASS` while build authorization is `NO-GO`
   and the verdict is `REVISE`.
3. A79 harness validation says `pass` in
   `dry_run_no_blender_no_candidate` mode; material validation says `PASS` for
   a standalone swatch only; geometry reports say `PASS` at pure-Python,
   pre-Blender, pre-render scope.
4. Early A79 projection proof measured the outer perimeter but omitted the
   bridge/whole-pocket envelope. Later evidence found bridge extrema and an
   inner field silently rewritten by as much as 3.593 mm / 71.09 degrees.
5. Sparse roots, aggregate 80% coverage, union-only brown masks, permissive
   angle/aspect thresholds, and omitted crossing/contiguous-failure checks
   could pass floating roots or an oversized planar shield. More validators
   were then added serially instead of asking the cheaper visual question.

Every verdict must therefore be named by stage, for example
`PURE_TOPOLOGY_PASS_RENDER_NOT_AUTHORIZED`; bare `PASS` is banned for partial
evidence.

## Duplicate and off-critical work

- Four approximately 9.47 MB frozen-geometry JSON files form two byte-identical
  pairs. Digest links would remove two duplicate payloads (about 18.9 MB) and
  one replay per state.
- The visual-hull branch produced 131 files / about 42 MB, including 120
  derived frame artifacts, only to conclude that calibrated hull recovery was
  invalid (12.54% median and 38.35% maximum antipodal disagreement). That
  scalar contradiction should have stopped the branch before mass image
  emission.
- A78 emitted eight substantive files for up to 6m 01s after its reset. Late
  workers were not cancelled or redirected promptly.
- A79 ran at least 13 visible workstreams: source/reference, registration,
  visual hull, interface, density, variants A/B/C, selector, render harness,
  build/scaffold/audit, material, rear-panel, and diagnostic work. Early
  reference/interface/harness concurrency was useful; the later output
  exceeded coordinator review capacity and did not shorten the serial build
  authorization choke.
- Material swatches and the separate rear-panel proposal were explicitly
  noncritical. Parallel execution did not prove they added wall time, but they
  added review/integration load before the macro geometry gate passed.

## Top five bottlenecks

| Rank | Bottleneck | Measured effect |
| ---: | --- | --- |
| 1 | Pre-render contract and authorization expansion | A79 waited 94m 30s for pixels that took 29.115s to render; 45m 36s elapsed after the process itself said “pixels next.” |
| 2 | Wrong/repeated representation boundaries | A77 + A78 spent 58m 58s canonically and produced zero candidate renders; A78 repeated a receiver family already rejected by history. |
| 3 | Partial `PASS` labels and incomplete proof subjects | Reopen, topology, harness, material, density, and pure-geometry passes did not imply render eligibility or likeness; later gates invalidated the branch. |
| 4 | Parallel fan-out beyond integration capacity | Roughly 399 A79 files, 23 Markdown files, 25 Python files, and 13+ workstreams accumulated before one diagnostic sheet. |
| 5 | Duplicate/unstopped off-critical work | Two duplicate 9.47 MB pairs, a 131-file hull branch, eight late A78 outputs, plus deferred material and rear-panel work before macro pixels passed. |

## Optimized next-cycle critical path

Use this strict sequence:

1. **0--5 minutes:** freeze the source hash and replacement boundary; check
   finite geometry, gross topology, non-overwrite output, and obvious source
   safety only.
2. **By 10 minutes:** render a disposable 320--512 px front plus worst
   profile/three-quarter in whole-plush context. This is diagnostic and may be
   non-promotable.
3. **One-minute self-veto:** reject bald coverage, helmet/curtain/card reads,
   gaps, and broken layer order directly from pixels. Stop all branches on a
   categorical failure.
4. **Only after a pixel survivor:** run exact root/crossing/thickness checks,
   clean reopen, semantic masks, and the full five-view packet in parallel.
5. **Only after that packet passes:** run independent blind review, material
   coupons, detailed density work, and promotion validation.

Prepare the generic disposable preview consumer and two-view compositor in
parallel with reference extraction, so geometry-to-pixel latency stays near
the observed 80 seconds. Keep one primary representation active; a backup
variant may run only until the first is renderable, then pause it. Cancel or
redirect workers immediately when a representation resets. Partial checks
publish scoped verdicts and link immutable payloads by hash instead of copying
them.

This retains the valuable safety and fidelity gates while moving silhouette,
coverage, curtain/card read, and layer order back to the only evidence that can
answer them: rendered pixels.
