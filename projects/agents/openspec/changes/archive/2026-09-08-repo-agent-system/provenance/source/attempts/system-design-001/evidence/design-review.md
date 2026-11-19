# Final design review

## Verdict

**Proceed.** The current candidate satisfies the design objective without the
authority duplication, synchronization cost, or bootstrap dependency of a
central agent brain.

## Strongest rejection case

The strongest objection to a distributed architecture is that owner-local
facts leave an agent repeatedly traversing prose and build graphs. The
strongest objection to the user's public-information boundary is that even a
non-secret operational fact can make an attack easier.

The candidate resolves the first objection with a thin shared contract,
bounded deterministic projections, one-hop root navigation, exact provenance,
and a roadmap for offline context and impact planning. It does not resolve it
by copying facts into another mutable store.

The second objection does not justify a vague confidentiality class in a
public-source repository. Origin-based secrecy is neither dependable nor
agent-legible. The design instead separates confidentiality from action
admission: credentials, other secrets, and personal information are never
disclosed, while authorization, environment binding, effects, and external
state continue to constrain what an agent may do with public operational
facts.

## Alternatives compared

- A hand-maintained central manifest improves apparent discoverability but
  creates duplicate fact authorities and stale safety decisions. Rejected.
- A mandatory daemon, database, or vector store can accelerate broad queries
  but adds a fallible bootstrap dependency before measured need exists.
  Deferred and not required by the architecture.
- A mega-skill reduces routing decisions but couples unrelated domains and
  makes regression ownership unclear. Rejected.
- Physical reorganization makes some paths look coherent while breaking
  natural subsystem ownership and does not solve provenance. Rejected.
- Thin contracts plus owner-local facts and bounded derived views preserve
  authority, degrade honestly, and permit incremental implementation.
  Selected.

## Incorporated adversarial findings

The final architecture separates `ValidationSet`, `EvidenceEvaluation`, and
goal-owned `EvidenceAssertion`; includes plans, admission, authority, review,
release, deployment, outcomes, and learning on the provenance spine; defines
history writes and distinct ownership axes; specifies the temporal loop,
runtime identity, failure and retry taxonomy, concurrency, hierarchical
budgets, and regression-first learning.

The final roadmap moves recoverable goals and untrusted runtime containment
ahead of higher-level joins, defines a closed registered inventory universe,
completes deterministic impact and broad-check workflows, adds release-ref
planning and receipts, and binds resource claims to measured baselines and
ceilings.

## Independent final checks

- Zero-context review: **PASS**, including root discovery, current/target/future
  separation, authority ownership, goal navigation, rendered targets, and
  revision-pinned evidence.
- Link review: **PASS**, including all five rendered pages, internal site
  targets, 16 pinned snapshot references, live authority links, source links,
  and reference definitions.
- Public-information review: **PASS**. No maintained document retains a
  “sensitive operational detail” class, and ignored secret scratch is not
  conflated with committed source.

The closed audit evidence remains immutable. Superseded proposals inside that
historical attempt are evidence of the audit process, not maintained policy.
