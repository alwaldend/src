# Commit messages and trailers

Read when preparing the aggregate commit message for `prepare`. The delivery
tool enforces these rules in `withCommitDisclaimer` and
`validateCommitFooter`, so an inaccurate message is rejected rather than
silently normalized.

Use concise, unprefixed Git subjects. A subject matching a Conventional
Commit prefix such as `type: ` or `type(scope)!: ` is refused as an
artificial prefix. Separate the subject from the body with a blank line, and
leave a blank line before the trailer footer. The footer holds only
`Token: value` Git trailers; any other footer line is rejected.

The trailer order is:

1. Optional `Learning-Proposal`.
2. `OpenSpec-Change` with the change path for OpenSpec-linked delivery
   commits.
3. Paired `Goal-Ref` and `Attempt-ID`, retained only for legacy goal-linked
   commits. Both are required together; supplying either one alone is
   rejected.
4. The generated `LLM-disclaimer` trailer, always final.

The tool appends the `LLM-disclaimer` commit trailer itself and rejects a
message whose final disclaimer line is incorrect, so callers normally need
not add it. The exact marker text is in
[the Forgejo compatibility workflow](forgejo.md).
