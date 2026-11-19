# Preserved goal provenance

This directory retains the original maintained record bytes from source commit `550d7e79b1f5fdbc2b6017b75178471d6914082f`. The [manifest](manifest.json) binds every original repository path to its snapshot path, SHA-256 digest, and size. The OpenSpec proposal, design, tasks, and migration metadata own the imported work; source files here are immutable evidence, including old workflow instructions and statuses.

Internal record-relative links keep their original layout. Links that depended on the old repository location are resolved below without rewriting the historical bytes. Unavailable task scratch remains unavailable; the migration does not invent artifacts or rerun past operations. Bare code-formatted paths and old external URLs remain historical locators. Current owner paths may have changed since the source candidate and are navigation only, not proof of historical bytes.

## Historical relative-link resolution

| Source document | Original target | Navigation at migration |
| --- | --- | --- |
| `projects/agents/goals/README.md` | `../docs/current-state.md` | [current owner path](../../../../../docs/current-state.md) |
| `projects/agents/goals/repo-agent-system/attempts/system-audit-001/evidence/delivery.md` | `../../../tools/repo_delivery/main/go/receipt.go` | Unavailable at migration; the historical target was not a tracked retained artifact. |
| `projects/agents/goals/repo-agent-system/attempts/system-audit-001/evidence/delivery.md` | `../../../tools/repo_delivery/main/go/delivery.go` | Unavailable at migration; the historical target was not a tracked retained artifact. |
| `projects/agents/goals/repo-agent-system/attempts/system-audit-001/evidence/delivery.md` | `../../../tools/repo_delivery/main/go/forge.go` | Unavailable at migration; the historical target was not a tracked retained artifact. |
| `projects/agents/goals/repo-agent-system/attempts/system-audit-001/evidence/delivery.md` | `../../../tools/versioning/skills/versioning/SKILL.md` | Unavailable at migration; the historical target was not a tracked retained artifact. |
| `projects/agents/goals/repo-agent-system/attempts/system-audit-001/evidence/delivery.md` | `../../../projects/goal/skills/goal/references/lifecycle-and-evidence.md` | Unavailable at migration; the historical target was not a tracked retained artifact. |
| `projects/agents/goals/repo-agent-system/attempts/system-audit-001/evidence/delivery.md` | `../../../projects/agents/skills/full-repo-check/scripts/run_full_repo_check.go` | Unavailable at migration; the historical target was not a tracked retained artifact. |
| `projects/agents/goals/repo-agent-system/attempts/system-audit-001/evidence/delivery.md` | `../../../projects/agents/skills/repo-delivery/SKILL.md` | Unavailable at migration; the historical target was not a tracked retained artifact. |
| `projects/agents/goals/repo-agent-system/attempts/system-audit-001/evidence/delivery.md` | `../../../tools/repo_delivery/main/go/command.go` | Unavailable at migration; the historical target was not a tracked retained artifact. |
| `projects/agents/goals/repo-agent-system/attempts/system-audit-001/evidence/delivery.md` | `../../../tools/repo_delivery/README.md` | Unavailable at migration; the historical target was not a tracked retained artifact. |
| `projects/agents/goals/repo-agent-system/attempts/system-audit-001/evidence/delivery.md` | `../../../tools/bazelrc/project.bazelrc` | Unavailable at migration; the historical target was not a tracked retained artifact. |
| `projects/agents/goals/repo-agent-system/attempts/system-audit-001/evidence/delivery.md` | `../../../tools/repo_delivery/main/go/git.go` | Unavailable at migration; the historical target was not a tracked retained artifact. |
| `projects/agents/goals/repo-agent-system/attempts/system-audit-001/evidence/delivery.md` | `../../../tools/versioning/cmd/versioning/command.go` | Unavailable at migration; the historical target was not a tracked retained artifact. |
| `projects/agents/goals/repo-agent-system/attempts/system-audit-001/evidence/delivery.md` | `../../../tools/git/main/bzl/al_git_repo.bzl` | Unavailable at migration; the historical target was not a tracked retained artifact. |
| `projects/agents/goals/repo-agent-system/attempts/system-audit-001/evidence/delivery.md` | `../../../tools/release/main/bzl/al_release.bzl` | Unavailable at migration; the historical target was not a tracked retained artifact. |
| `projects/agents/goals/repo-agent-system/attempts/system-audit-001/evidence/delivery.md` | `../../../projects/ci_platform/releases/BUILD.bazel` | Unavailable at migration; the historical target was not a tracked retained artifact. |
| `projects/agents/goals/repo-agent-system/attempts/system-audit-001/evidence/delivery.md` | `../../../tools/release/main/go/generator.go` | Unavailable at migration; the historical target was not a tracked retained artifact. |
| `projects/agents/goals/repo-agent-system/attempts/system-audit-001/evidence/delivery.md` | `../../../tools/release/README.md` | Unavailable at migration; the historical target was not a tracked retained artifact. |
| `projects/agents/goals/repo-agent-system/attempts/system-audit-001/evidence/delivery.md` | `../../../tools/release/main/go/deployer.go` | Unavailable at migration; the historical target was not a tracked retained artifact. |
| `projects/agents/goals/repo-agent-system/attempts/system-audit-001/evidence/delivery.md` | `../../../AGENTS.md` | Unavailable at migration; the historical target was not a tracked retained artifact. |
| `projects/agents/goals/repo-agent-system/attempts/system-audit-001/evidence/delivery.md` | `../../../projects/agents/skills/git-rebase-remote/SKILL.md` | Unavailable at migration; the historical target was not a tracked retained artifact. |
| `projects/agents/goals/repo-agent-system/attempts/system-audit-001/evidence/delivery.md` | `../../../projects/agents/skills/git-rebase-remote/BUILD.bazel` | Unavailable at migration; the historical target was not a tracked retained artifact. |
| `projects/agents/goals/repo-agent-system/attempts/system-audit-001/evidence/delivery.md` | `../../../projects/goal/skills/goal/references/record-format-v1alpha1.md` | Unavailable at migration; the historical target was not a tracked retained artifact. |
| `projects/agents/goals/repo-agent-system/attempts/system-audit-001/evidence/delivery.md` | `../../../tools/repo_delivery/main/go/review.go` | Unavailable at migration; the historical target was not a tracked retained artifact. |
| `projects/agents/goals/repo-agent-system/attempts/system-audit-001/evidence/delivery.md` | `../../../projects/agents/skills/bazel-nested-module/SKILL.md` | Unavailable at migration; the historical target was not a tracked retained artifact. |
| `projects/agents/goals/repo-agent-system/attempts/system-audit-001/evidence/delivery.md` | `../../../projects/agents/skills/repo-bazel/SKILL.md` | Unavailable at migration; the historical target was not a tracked retained artifact. |
| `projects/agents/goals/repo-agent-system/attempts/system-audit-001/evidence/delivery.md` | `../../../tools/git_hooks/main/go/precommit.sh` | Unavailable at migration; the historical target was not a tracked retained artifact. |
