## Context

Keep submission pages on the main site after their source relocation.

## Goals / Non-Goals

Update the existing submission package dependency to //users/simeonwarren/leetcode:leetcode_submissions. Retain the section and generated pages.

## Decisions

The user explicitly requires site inclusion from the user-owned source, overriding the generic user-tree production restriction for this target.

## Risks / Trade-offs

The full site builds with 1,928 submission pages plus the section index; a representative page retains its submission link and code. Exact candidate validation remains required before publication.
