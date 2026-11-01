## 1. SSH wrapper

- [x] 1.1 Add the AL-backed `ssh_binary`, default target, and documentation;
      verify the target builds and the host client's `-V` and `-G` output.

## 2. Integration verification

- [x] 2.1 Verify argument forwarding using offline SSH configuration output and
      run `@rules_al//cmd/al:al:command_test` for AL injection lifecycle coverage.

## Validation cases

Before implementation, check these possible failures: selecting a bundled
client instead of host `PATH`, losing argument boundaries, failing to forward
AL configuration or runtime data, and bypassing injector cleanup. Use offline
client inspection and AL's existing end-to-end injector lifecycle test; no live
authentication or SSH connections are needed.

The optional design artifact is omitted because this macro only composes the
existing AL wrapper and introduces no new lifecycle or authentication design.

## Evidence

On 2026-10-25, `al bazel run //tools/ssh -- -V` reported the host
OpenSSH client. An offline `-G -F /dev/null` invocation preserved the requested
user, port, and identity path containing spaces. The AL command test passed
its injector success, startup failure, and cancellation scenarios. Local
artifacts are in ignored `out/ssh/`.

Package lint and repository quality are mandatory publication gates through
`al-delivery`; Git and delivery receipts own commit, push, and PR state.
