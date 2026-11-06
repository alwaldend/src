## Evidence

The pre-implementation fixture failed after eight seconds because the Terraform
descendant retained the caller's output pipe. The fixture starts the actual
rental executable and actual packaged launcher, then signals only the rental
process while its Terraform fixture waits for cancellation.

After graceful forwarding at both boundaries, the same fixture passes: Terraform
records the interrupt, release removes the partial contract, and the command
writes a receipt with `cleanup_attempted` and `state_empty` both true. This is an
offline process lifecycle check; it does not establish live provider acceptance.

All 16 standalone-module tests and both owning OpenSpec validators pass. The
initial standalone invocation used a noncanonical TMPDIR and exposed existing
literal-path assertions; rerunning with the canonical task scratch path passes.
No source change was needed for that invocation error.

Graceful shutdown allows ten seconds at the Terraform boundary and twenty at
the rental caller. Hard termination or failed cleanup still requires inspection
and the rental test state's owning recovery command.
