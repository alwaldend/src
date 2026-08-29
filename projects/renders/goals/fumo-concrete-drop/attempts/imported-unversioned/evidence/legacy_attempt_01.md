# Attempt 1 — neutral rigid drop scaffold

[Back to goal](../README.md)

## Plan

The frozen hypothesis, inputs, parameters, and rejection conditions are in
[the current-attempt record](../current_attempt.md).

## Work, evidence, and decision

Blender MCP created one candidate from factory settings.  It contains the
`FUMO`, `SET`, `LIGHTS`, and `CAMERAS` collections, a
`FUMO_Drop_Interface`, a neutral `Fumo_Rig_PLACEHOLDER`, a passive concrete
floor, and an active box collider.  The protected blend hashes did not change.

Candidate:
`86630e599525e40663ad01e4bd8f4c5c5f12e9cb127740440a7bbe501b77d292`.
It descends `.84534 m`, contacts at frame `22`, penetrates the sampled floor by
at most `.000687 m`, and has zero measured late motion.

The [contact sheet](../../../../../out/fumo_concrete_drop_scaffold/attempt_01/contact_sheet.png)
is nevertheless an absolute visual failure: extreme clipping makes the floor
and proxy nearly white, the settled proxy is too small, and the warning text
is absent.  Decision: **reject Attempt 1 as an inspectable scaffold while
retaining its mechanics and interface settings for Attempt 2**.
