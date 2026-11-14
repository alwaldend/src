# ComfyUI installation acceptance

Run `bazel test //projects/ansible_collection/test/comfyui_tools:installation_test`.
The packaged role runs through local Ansible against an isolated application
tree, a tiny HTTP model/readiness fixture and a local Git repository using the
host Git executable (a repository prerequisite). The restart command only
records fixture events; it never calls the host service manager.

The fixture verifies a non-default workspace, checkpoint/LoRA/ControlNet paths,
exact Git revision, Python/browser source and App preservation, no restart on
unchanged application, a changed-node restart while skipping workflows, busy
queue refusal, checksum refusal and path traversal refusal. It emits
`installation.json` in Bazel's undeclared outputs. GPU loading, real Python
dependency installation, Vulkan discovery and SSH authentication remain outside
this fixture and require separately authorized deployment.
