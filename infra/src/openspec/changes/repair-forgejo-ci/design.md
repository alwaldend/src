## Context

The supplied run built commit `8b91edb37de21007becef5a37d41ffef2371f102`.
Release generation could not traverse a shallow checkout; Android AddToZip
reported `zip: command not found`; five diagram projections differed; the
secret scanner reported five unmatched findings. External downloads failed
with geographic rejection, forbidden responses, DNS failures, or timeouts.
A read-only HashiCorp HEAD probe on 2026-11-17 returned 404 with
`x-amzn-waf-reason: geo`. Changing valid dependency versions does not fix that.

## Decisions

Use full checkout history, preserving the event SHA, Git LFS and disabled
credential persistence. Keep all executable work in the existing CI command.
Keep runner changes in Ansible and reuse an existing persistent HTTP proxy,
excluding local destinations and `.alwaldend.com`. Never copy the temporary
loopback proxy injected into this agent session.

Regenerate image projections through the existing update target. Review exact
historical source coordinates before adding allowlist entries; add inline
comments for current synthetic fixtures. Do not suppress whole files or detectors.

## Validation and continuation

Raw supplied-log analysis and command output remain in ignored
`out/forgejo-ci/`. Validate offline targets before any authorized deployment.
Record actual commands, candidate, proxy source and live limitations here or
in linked owner-local evidence. Publication is authorized; merge is not.
