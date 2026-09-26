# Runner acceptance

These cases were written before runner implementation. The executable acceptance harness
launches the same Bazel-packaged runner as consumers; no VM lifecycle mocks
are permitted. Each invocation retains `result.json` and selected diagnostics.

| Case         | Stimulus                                                                   | Required observation                                                                                                            |
| ------------ | -------------------------------------------------------------------------- | ------------------------------------------------------------------------------------------------------------------------------- |
| Smoke        | Fresh guest, package installation, file convergence, unchanged run, reboot | SSH and package access succeed; second convergence changes nothing; marker survives; image digest unchanged                     |
| Failure      | Fail a named verification assertion after convergence                      | Nonzero exit; failed phase and assertion recorded; VM gone; private state removed                                               |
| Cancellation | Send SIGTERM while verification waits                                      | Nonzero exit within grace deadline; interrupted phase recorded; VM gone; private state removed                                  |
| Concurrent   | Start two smoke runs together                                              | Different run IDs, ports, disks, keys, and QMP endpoints; both complete independently                                           |
| Repeat       | Execute smoke twice with identical inputs                                  | Different run IDs; both actually execute rather than returning cached results                                                   |
| Prerequisite | Select an invalid accelerator or invalid image                             | Actionable preflight failure before convergence; no surviving VM                                                                |
| Isolation    | Set sentinel Vault, AWS, SSH-agent, and Ansible environment values         | Guest and controller see none; retained files contain no sentinel or generated private credentials                              |
| Recovery     | Kill supervisor abruptly after VM starts, then invoke scoped destroy twice | Exact owned process is removed; unrelated process survives; second destroy succeeds                                             |
| Mapping      | Stage an executable fixture at a renamed collection destination            | Guest receives those exact executable bytes; candidate mapping hashes identify the fixture                                      |
| Evidence     | Inspect each result                                                        | Tool/image/input identities, accelerator, phases, package versions, assertion summary, rerun label, and cleanup outcome present |

QEMU uses loopback forwarding and disposable file-backed disks. Tests must
never create a host bridge, install host packages, use a production inventory,
or read developer credentials. A failed cleanup fails acceptance.
