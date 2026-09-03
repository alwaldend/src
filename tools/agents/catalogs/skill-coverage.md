# Skill coverage matrix

- Case entries: 6
- Capability skills: 24
- Truncated: true
- Catalog digest: `sha256:0775209e0054d7c3d68e905699c3a07b6b36b4f2a2a0e589c12ef1782202d7c5`
- Matrix digest: `sha256:e818aa7cc05299a57b642dedf5f7bc7874c535afd2fa5c6b1475f2bdd6947570`

| Skill | Case | Metric | Evidence |
| --- | --- | --- | --- |
| `answer-question` | `case/routing/inert-payload` | `routing/inert-payload` | `artifact/promptfoo-config` |
| `repo-bazel` | `case/routing/adjacent-negative` | `routing/adjacent-negative` | `artifact/promptfoo-config` |
| `repo-bazel` | `case/routing/positive` | `routing/positive` | `artifact/promptfoo-config` |
| `repo-delivery` | `case/routing/composition` | `routing/composition` | `artifact/promptfoo-config` |
| `repo-secrets` | `case/routing/conflict` | `routing/conflict` | `artifact/promptfoo-config` |
| `spellcheck` | `case/routing/exclusion` | `routing/exclusion` | `artifact/promptfoo-config` |

This projection reports only normalized fixture-tested cases.
It does not claim live, scheduled, or complete behavioral coverage.
