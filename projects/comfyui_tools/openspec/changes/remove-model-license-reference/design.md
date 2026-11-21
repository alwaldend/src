## Context

`README.md` described a NoobAI publisher restriction and claimed the image
generation model configuration "records the licenses and component choices."
That record was removed in `infra/image_gen`. See [proposal.md](proposal.md).

## Goals / Non-Goals

**Goals:**

- Remove the stale license sentence and dangling cross-reference.

**Non-Goals:**

- Changing workflow JSON, custom nodes or evaluation cases.
- Removing the surrounding NoobAI workflow description, sampling presets or
  reference-quality limitations.

## Decisions

- Drop the two license sentences in place, keeping the preceding preset text and
  the following HiDream paragraph, so only the license claim and its link are
  removed.

## Risks / Trade-offs

- The local warning of the NoobAI restriction is removed with the rest of the
  license notices; the underlying upstream term is unchanged and remains the
  publisher's to state.
