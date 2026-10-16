# Reimu Fumo process review and research

Observed 2026-09-21 by `/root/fumo_process_review`. Task: rebase, review the
process, research improvements, and continue the existing Fumo. This is
task-local evidence and advice, not canonical policy or an acceptance record.
The coordinator owns Git, native Blender, the goal, and the final decision.

**Coordinator verdict: REVISE.** Continue with one bounded hair B fit-only causal
correction plus diagnostic facial footprints. The hair diagnosis inspected
canonical turn frames 12/21, which show real long lobes and bow-tail occlusion;
inventing darts or simultaneously changing patterns would confound the fit
test. B can establish fit only. No fiber, rigging,
detail work, or further body A/B parameter sweep. Preserve head109 and arm101;
neither is an accepted whole-character foundation.

## Scope and evidence

Read-only review covered the asset/process instructions, saved 115 state,
historical process reviews 018/047/104, results 107 and 112–114, and 115's builder.
The reviewer directly inspected canonical front and physical side references,
hair 115A close front/side/three-quarter, and body 114B front/side pixels. Having
read implementation and history, this reviewer is not an implementation-blind
acceptance reviewer. Historical results were consumed as recorded evidence;
old model outputs were not independently reproduced.

All local paths below are relative to `out/reimu_fumo_finish/desktop_astra/`:

| Evidence                                         | Exact identity                                                                                                        |
| ------------------------------------------------ | --------------------------------------------------------------------------------------------------------------------- |
| Hair A candidate, recorded by clean-open receipt | `hair_parallel_115_a.blend`, SHA-256 `46020796b7bb913e2acd1dd895e70eb1a1ff8efad1d0293d2d764801b05b335d`               |
| Hair A render receipt, rehashed in this review   | `hair_parallel_115_a_fast/receipt.json`, SHA-256 `3772803ba211a4814b3e543fc4b678f1acc6d9c2affad7a233ebdc5160f78f6d`   |
| Hair A builder, read and rehashed                | `hair_parallel_115_build.py`, SHA-256 `afd0ea6b7ecff5e13d4a9afd26551ce72257d21c32055e5696a3835d1879cd00`              |
| Body B candidate, recorded by clean-open receipt | `body_parallel_114_candidate_b.blend`, SHA-256 `c0c41fd23e4e5fc862c756d19844046329ba6704e618f11dd421da741e7bf048`     |
| Body B render receipt, rehashed                  | `body_parallel_114_review_b/receipt.json`, SHA-256 `38f51cdf7f176e22695333b74c3fae86f3a70fc8e055abf54bef085ed67c1048` |
| Prior macro strategy review, rehashed            | `process_104_ground_up_review.md`, SHA-256 `242bf3acfd67b6cff637b35739bd6e32daeac872f4580e4b35023be8dc95bf4e`         |

Receipts identify Blender 5.2.1 LTS build `9e2066aef7ef` and review contract
`4835f1595995db408567044849ff8f2f19717b9ce1a6492fc85de34755ac7be4`.
Candidate hashes above are receipt bindings, not newly repeated candidate-byte
audits. Source inspection before this write found branch
`t3code/continue-fumo-desktop-use`, HEAD `c7601f0fc80e0a94b585910459639ffccbdbdbd4`,
in its linked worktree; this observation does not claim the rebase completed.

## Findings that change the next modeling step

- **FUMO-EVALUATED-FIT — live pixels and inspected source.** Hair 115A exposes
  broad head surfaces through the crown/fringe. Its builder fits selected cage
  rows, hardcodes four upper arches, then subdivides paired surfaces without an
  evaluated clearance check. The historical 047 audit recorded a related
  cage-versus-surface crossing. Proposed correction: constrain the fitted
  crown/fringe after subdivision, leave free hems outside that constraint, and
  check evaluated vertices plus face samples before inspecting two views.
  Measure success by disappearance of the exposed arcs without displaced
  canonical landmarks, detached roots, or a new silhouette regression.
- **FUMO-PROCEDURAL-LOCKIN — direct 115 pixels; historical 112–114 records.**
  Shared adjacency and new quads solved some technical problems but did not
  remove the helmet, curtain, and blade appearance. B isolates fit; diagnostic
  contact improvement cannot establish that construction is resolved. Cape or
  blade construction remains an independent rejection risk, judged against
  the canonical lobes and occlusion rather than an invented replacement pattern.
  Retention still requires both viable fit and construction. Different topology
  alone is not a successful representation change; do not assemble rejected
  hair or hide it with texture.
- **FUMO-PATTERN-CAUSE — direct body B pixels.** The skirt reads as a tent over
  an intersecting torso, with a mechanical hem. Its A/B budget is exhausted.
  A future body strategy should originate fullness at a gathered waist over a
  seated core, then establish foot contacts. Another radial cone with waves
  around its edge would not test that construction hypothesis.
- **FUMO-MACRO-FOUNDATION — historical evidence.** Review 104 already called for
  a macro reset; 105/106 then failed, and 107 established softness without the
  required sleeve shape. A clean rebuild or generic mesh is therefore not
  itself a breakthrough. Return a viable subsystem promptly to whole-character
  context before claiming that it solves the user's visual complaint.
- **FUMO-RESUME-DRIFT — live bounded measurement.** `CURRENT.md` has 171 lines /
  10,975 bytes and stale future instructions; the 61-line save checkpoint
  explicitly supersedes them. Keep a short current index with actual bytes,
  verdict, protected fallback, and next action; link historical explanations.
- **FUMO-FEEDBACK-INVERSION — live logs and historical review.** Exact receipts
  have repeatedly preceded obvious image rejection. Current rendering is
  already cheap: hair's five views took about 12 seconds and body's eight about
  24 seconds. Improve visual diagnosis before adding tooling to optimize these
  small costs. Keep clean reopen and independent retention review.

## Primary-source research and limits

Blender's [Shrinkwrap manual](https://docs.blender.org/manual/en/4.2/modeling/modifiers/deform/shrinkwrap.html)
documents vertex-group influence and offsets along the target's smooth normal.
This supports selective fitting after subdivision while leaving hanging cloth
free. Its inside/outside mode is only a crude collision check around sharp
corners; it cannot replace inspection of evaluated surfaces. This is tool
capability evidence, not proof that B will look correct.

Choly Knight's [Bat & Cat Girl tutorial](https://cholyknight.com/wp-content/uploads/2019/10/bat-cat-girl-doll-plush-sewing-pattern.pdf),
pp. 6, 8–10, constructs a padded face with darts and attaches overlapping bangs.
The [Jack & Sally tutorial](https://cholyknight.com/wp-content/uploads/2022/09/Jack-Sally-Doll-Plush-Sewing-Pattern.pdf),
pp. 27–29, uses doubled, darted hair fitted around the head. **Inference:**
perimeter-controlled padded volume and shaped cloth panels provide better
construction causes than uniform inflation. These different dolls do not
establish Reimu's exact pattern or proportions.

Lisa Press's [gathered-skirt tutorial](https://www.phoebeandegg.com/blog/2014/9/16/doll-dressmaking-seriesdress-with-a-bodice-and-gathered-skirt)
starts with a rectangle gathered at its upper edge to fit the bodice. **Inference:**
waist attachment and excess fabric should cause skirt fullness; the controlling
Fumo images still decide its seated silhouette and contacts.

The strongest alternative is native cloth simulation or
[Cloth Filter](https://docs.blender.org/manual/en/latest/sculpt_paint/sculpting/tools/cloth_filter.html),
which supports masked/Face Set pinning and inflation/expansion. That may improve
softness, but 107 directly disproves the inference that soft folds imply correct
sleeve construction. Defer simulation until placement and pattern shape are
credible. Doing nothing preserves bytes but does not satisfy the continuation;
a bounded B is reversible and tests a specific diagnosed cause.

## This turn's additional ergonomics evidence

The coordinator reports three observed costs; this reviewer did not reproduce
them: an unbounded worktree list and broad tool metadata caused output
truncation; installed `bazel_agent` lacks the assumed `tool` subcommand; and a
cold delivery build encountered contention with the Cordis server's Bazel work.
The coordinator's `jcmd` evidence identified a thread actively decompressing XZ
data, rather than a network hang. The original 180-second timeout interrupted
the cold dependency build; a causally justified 900-second retry is running at
this update. This is observed setup progress, not a completed-build claim.
Use bounded current-worktree facts and selected tool metadata, consult the
installed runner's actual supported interface once, and reuse that result.
Keep dependent Bazel work sequential and distinguish cold setup from failure;
preserve the runner's locking behavior rather than bypassing it. Measure
improvement by zero repeated interface discovery/truncation and no duplicated
cold setup. These observations do not authorize host or shared-tool changes.

The next decision remains visual: B must improve fit in front and side without
a new critical regression. This fit-only experiment cannot claim to resolve
cape or blade construction; retention requires a separate favorable judgment
of those forms in the same pixels. Diagnostic facial footprints
may check placement; they do not advance the material tier or pass a criterion.
No full model is accepted by this review.
