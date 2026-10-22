## 1. Assess maintained work

- [x] 1.1 Inventory active changes at the reviewed candidate and retain the original checked/unchecked counts in task-local evidence.
- [x] 1.2 Verify each archive candidate against its artifacts, implementation, and retained evidence; record per-change readiness and per-delta conflict decisions.

## 2. Reconcile and archive

- [x] 2.1 Merge implemented requirements into current owner-local baselines and archive eligible changes, verifying that superseded contracts remain historical and outstanding work stays active.
- [x] 2.2 Repair maintained references and include affected owner archives in structural validation; verify archive/source membership against the final inventory.
- [x] 2.3 Create valid, explicitly deferred owner-local plans for findings 1–4 and 6, verifying that implementation tasks remain unchecked and no proposed delta enters a current baseline.

## 3. Verify the result

- [x] 3.1 Run preflight strict active/archive OpenSpec validation, formatting, quality, and affected semantic checks; inspect the specifications and archive audit and select the final exact-candidate delivery gates.
- [x] 3.2 Reconcile final counts and preserved active changes against the original inventory, record validation limits, and archive this maintenance record after source acceptance. The delivery receipt owns final exact-candidate validation and publication.

Preflight evidence: aggregate OpenSpec 85/85, repository quality 26/26,
semantic lint of the changed validation packages, and the parent delivery
documentation build passed. Independent reviews verified source contracts and
archive integrity. The final validation plan repeats aggregate OpenSpec and
quality, lints all affected packages, covers nested-module gaps and builds the
documentation consumer against the delivered commit.
