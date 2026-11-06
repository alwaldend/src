## Context

The current infrastructure skill recommends saving plans. Some component
targets enforce that contract, while ordinary apply targets support an
interactive prompt. See proposal.md for the requested change.

## Goals / Non-Goals

Apply the agreed agent workflow through policy and canonical skill guidance.
This change does not alter reusable Terraform rules, component command
contracts, live infrastructure, or other tasks' files.

## Decisions

- Keep the normative rule in root AGENTS.md and the operation procedure in
  repo-infra. The skill references policy instead of maintaining another owner.
- Use Terraform's existing sensitive-value redaction. Normal text in agent
  transcripts is intentional; additional output filtering is not required.
- Report saved-plan-only entry points as incompatible. Silently dropping
  their guard or choosing another target could broaden an authorized operation.
- Audit leftovers through filenames, file metadata, and ZIP member names;
  never display plan payloads. Preserve unrelated artifacts and report coverage.

## Risks / Trade-offs

- Unmarked secrets can appear in text output; the user accepts this risk and
  the existing sensitive declarations remain necessary.
- Some scoped deployments remain unavailable to agents until their owning
  command contract is separately updated.
- A bounded filename audit cannot prove absence of arbitrarily named plans;
  report its limits instead of treating no matches as exhaustive proof.
