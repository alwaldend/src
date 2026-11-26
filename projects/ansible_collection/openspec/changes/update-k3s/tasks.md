## 1. Upgrade K3s

- [x] 1.1 Verify the upstream stable release and pin the binary and installer source.
- [x] 1.2 Pass the force-restart option when the binary copy changes.

## 2. Verify the candidate

- [x] 2.1 Confirm the pinned artifacts, reported version, and restart behavior for changed and unchanged binaries.
- [x] 2.2 Pass the Flux and Harbor Ansible syntax checks and validate the OpenSpec change.

The downloaded binary matched the upstream amd64 checksum and reported
`v1.37.2+k3s1`. The source archive matched its Bazel integrity pin.
Two localhost Ansible runs passed `true` and `false` to the installer's
restart option on changed and unchanged binary copies, respectively.
Both deployment syntax tests passed in a workspace-local Bazel sandbox.
