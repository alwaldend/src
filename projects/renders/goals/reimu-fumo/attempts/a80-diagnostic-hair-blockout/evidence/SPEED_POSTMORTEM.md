# A80 loop-speed postmortem

Snapshot: 2026-08-31 16:10 +03. Filesystem times identify artifact
completion, not worker start. The formal start is the canonical attempt
creation time, `2026-08-31 15:00:44.756 +03`.

## Verdict

A80 fixed A79's gate-ordering and artifact-sprawl failure, but not its
pre-render latency. First whole-character pixels arrived in **50m 41.672s**,
versus A79's **94m 30s**: **43m 48.328s faster (46.4%)**. A80 nevertheless
missed the A79-derived 12-minute first-pixel budget by **38m 41.672s**.

The two A80 Blender cycles totaled only **33.901s**; measured rendering was
**33.708s**, or **0.98%** of the **56m 48.242s** formal-start-to-final-sheet
span. The remaining opportunity is before renderer invocation and in the
render-to-review handoff, not in image quality, source-integrity checks, or
Blender throughput.

## Measured timeline

| Elapsed | Artifact/event | Evidence |
| ---: | --- | --- |
| +0 | Canonical A80 attempt created | `attempt.yaml` |
| +5m 16.726s | Empty `variant_a/` directory | filesystem birth time |
| +49m 29.995s | Sculpted one-off builder complete | script mtime |
| +50m 41.672s | First A80 pixel, sculpted front | PNG mtime |
| +50m 48.003s | Sculpted controlling pair complete | side PNG mtime |
| +50m 54.959s | Sculpted third view and initial manifest complete | rear/manifest mtimes |
| +54m 09.147s | Sculpted categorical self-review | review mtime |
| +54m 38.517s | Constructed front pixel | PNG mtime |
| +54m 43.464s | Constructed controlling pair complete | 3Q PNG mtime |
| +54m 48.221s | Constructed third view and manifest complete | rear/manifest mtimes |
| +56m 41.010s | Constructed categorical self-review | review mtime |
| +56m 48.242s | Final rejected sheet manifest | manifest mtime |

The first populated variant directory appeared only at 15:49:43.947, so
**86.4%** of the plan-to-final span elapsed before that point. A broad
repository checkpoint also completed inside this interval: commit `f312981f`
was authored at 15:25:42, committed at 15:42:36, and changed 104 files
(+9,841/-177). The timestamps prove overlap, not that the commit caused the
delay.

## What improved from A79

| Measure | A79 | A80 | Result |
| --- | ---: | ---: | --- |
| First whole-character pixel | 94m 30s | 50m 41.672s | 43m 48.328s faster |
| Pre-pixel/output volume | ~399 files / 148 MiB | 22 files / 6.5 MiB total at review | major churn reduction |
| Python authored | ~13,000 lines | 947 lines | major setup reduction, still duplicated |
| Promotion-grade work before visual veto | extensive | none | correct gate order |
| Beauty rendering | 5 views, 29.115s | 6 views, 33.708s | still negligible |

A80 correctly retained fixed cameras, whole-plush context, pinned Blender,
protected hashes, exact six-object replacement boundaries, explicit
`NON_CANDIDATE` labels, and immediate rejection before mechanics, materials,
semantic masks, reopen, or promotion work. Those controls did not create the
delay.

## Evidence-backed speed changes

1. **Keep repository checkpointing off the first-pixel critical path.** The
   only large measured interval is the first 49m30s before the sculpted builder
   appeared; the 104-file checkpoint occurred inside it. Start the visual
   cycle after such coordination work, or let the already-bounded geometry
   worker reach pixels without waiting for it. No exact saving can be assigned
   because worker start and concurrency were not recorded.

2. **Reuse one existing diagnostic runner and compositor; vary only the
   geometry.** A80 authored two build/render scripts totaling **796 lines /
   30,141 bytes** and two compositors totaling **151 lines / 4,830 bytes**.
   Both reimplemented source/hash checks, the six-object hide boundary,
   Blender checks, render configuration, save/render loops, hashing, and
   manifests. A79 had already identified a reusable disposable preview as the
   next-cycle requirement. Reuse targets the measured pre-invocation gap while
   preserving every safety and pixel-quality control; the current evidence
   does not justify claiming a specific minute saving.

3. **Stop a categorical loser after the declared two-view pair.** Sculpted
   front+profile already showed the dome, cuboid profile, and beige exposure;
   constructed front+3Q already showed the helmet band and bald patch. Their
   rear renders cost **6.956s + 4.751s = 11.707s**, **34.7%** of all A80
   render time. Keep rear/opposite-side/full-view renders mandatory for a
   survivor; they were unnecessary to reject these two variants.

4. **Review the pair as soon as its second PNG closes.** The sculpted pair sat
   **3m 21.159s** before its self-review; the constructed pair sat
   **1m 57.556s**. Pixel composition was not the cause: the recorded sculpted
   board operation took **1.017s**. Use the already-rendered pair for the
   categorical decision, then compose the annotated archival sheet after the
   stop. This changes ordering, not review quality.

5. **Finish the handoff when both verdicts exist.** At the 16:10 snapshot,
   more than 13 minutes after the final rejected sheet, the canonical attempt
   was still `open`, with no evidence and `result.md` still reporting no
   result. Closing immediately does not accelerate Blender, but it removes
   measured dead time before the next bounded attempt can become canonical.

## Do not optimize

- Do not lower survivor/final render quality. The 420 px constructed pair was
  enough for gross vetoes, but subtle seams, contacts, one-pixel gaps, and
  likeness still require 512--640 px, registered comparisons, and the full
  fixed-view set.
- Do not remove hash, path, immutable-source, or save checks. Manifest-accounted
  non-render overhead across both Blender cycles was only **0.193s**.
- Do not replace the candidate-only reject sheets with fewer quality gates for
  a survivor. A80's sheets lacked direct reference/rung-003 rows and were not
  implementation-blind; that is acceptable only because the visible failures
  were categorical. A survivor still needs cached `reference | rung003 |
  candidate` comparisons, complete direct views, and blind review.

The next-cycle target should therefore remain the A79 budget: first pair by
12 minutes and rejection by 15 minutes. A80 proves the renderer can satisfy
its part in seconds; achieving the budget requires removing work before
invocation and acting on the pair immediately.
