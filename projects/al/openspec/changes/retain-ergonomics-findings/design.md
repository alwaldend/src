## Delivery sequence

The delivery skill owns the timing and completion condition for the ergonomics
review. The existing ergonomics skill owns evidence collection and assessment.
Apply the same procedure to both forge paths before final preparation. Follow-up
outside the task can use an existing owner record or a handoff in the current
change. A review with no actionable findings needs no additional artifact.

## Observed evidence

The NoobAI deployment reported two avoidable repeats: a raw Terraform output
call for a nullable rental ID and a ComfyUI JSON query that tried to read system
stats at the top level. The deployment verified both workarounds and completed.
The broader defect was procedural: the dispositions lived only in ignored
scratch. This change addresses that defect at the shared delivery step.

The image-generation component owns the two command examples. In its next
runbook edit, use the normal Terraform output mode for a nullable instance ID
and select the system-stats payload from `.data`. Verify that the rental check
works with no active contract and that the authenticated stats example returns
the GPU name on an active contract.
