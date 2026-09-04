# Ergonomics-report provenance correction

This record corrects the source identity of
`flatten-dose-response-018/evidence/session_ergonomics_review.md` without
rewriting that closed attempt.

## Reachable anchor

- Observed at: `2026-09-04T14:52:29Z`
- Truncation status: `partial`. The Git anchors below are complete. The
  unavailable-input section lists the known missing input categories, but not
  every constituent path or byte because the session inventories were never
  frozen.
- Goal tree: `733966a68d40eb97dc1e08b22b6d02125dc5f704`
- Closed-attempt tree: `0a954fc6f9fe759427927d0b4fa91371da2ec89b`
- Report blob: `8533ae01c89cdd9a38e8dfba7b8ab273922b520d`
- Report path:
  `projects/renders/goals/reimu-fumo-finish/attempts/flatten-dose-response-018/evidence/session_ergonomics_review.md`
- Report SHA-256:
  `7edd22cbd6788af5dd1c2cca772f409996688b4e780f4526ece36ca178cda73d`

These Git object identities are taken from delivered snapshot
`1b785167e610d12bee3ac0a4c9be821b445cf1ca`. The goal and closed attempt are
unchanged by this correction, so their trees and report blob remain reachable
at the same paths in the correcting commit even if one-commit publication
changes its commit identity. The goal tree contains resource version 58 in the
final blocked state.

## Unavailable historical inputs

The report was measured earlier in worktree `t3code-c76e14e4` while the branch
temporarily pointed at
`e47f0102b19bbc45da7cb804a832d0f5b40450b1` and goal resource version 57 was
uncommitted. That superseded commit is not a reachable branch anchor for a
fresh checkout, and the exact uncommitted resource-version-57 bytes were never
frozen.

The following inputs were also session-local and are not recoverable from the
repository snapshot:

- ignored `out/reimu_fumo_finish/` scratch bytes used for the 382-file,
  235,026,682-byte, and file-type counts;
- ignored agent scratch used for its 84,072,622-byte count; and
- the live collaboration inventory used for the agent-state counts;
- the exact selected-skill manifest and versions behind the 11-file,
  78,031-byte context count;
- the user-message and interrupt timeline behind the feedback-latency
  findings; and
- the prior closeout review and scratch inventory behind the seven-attempt,
  approximately 142 MiB comparison.

Consequently, those counts, comparisons, and latency findings are historical
observations, not independently recomputable claims about commit `1b785167`.
They support the direction of the ergonomics recommendations but are not
acceptance evidence for the Reimu asset. The durable conclusions remain
unchanged: no improved model was retained, no criterion passed, and execution
is blocked pending a qualified writer or genuinely different proven authoring
route.
