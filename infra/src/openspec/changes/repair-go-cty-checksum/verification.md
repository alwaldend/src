The original broad semantic lint command succeeds after generating one archive
checksum line for the already selected go-cty v1.13.1. The Go manifest and all
dependency versions remain unchanged. The pinned Go downloader used the
configured sum.golang.org checksum database with root-module workspace mode
disabled for this invocation only.

The original failure and successful lint logs belong to PR 211's task evidence.
For standalone delivery, pinned Go regenerated the same single-line checksum
change directly on master. Repository quality, shared-tool semantic lint and
the existing HCL consumer build run against the prepared candidate. Logs and
local review evidence live under out/checksum-extraction; AL receipts own the
exact commit, remote branch and publication state.
