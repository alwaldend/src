# Bazel session cleanup

The repository's `.codex/config.toml` queues asynchronous Bazel expunge when a
Codex session ends. Review and trust the hook with `/hooks` in Codex; project
hooks do not run until trusted. A completed assistant turn does not trigger
cleanup.

If `/hooks` is empty in a linked Git worktree, check the primary checkout's
`.codex/config.toml`. In Codex CLI 0.153.4, observed on 2026-09-06, `config/read`
returned the primary checkout's configuration even though its project-layer
metadata named the linked worktree. That primary configuration had no hooks,
so `hooks/list` returned an empty list. The exact worktree TOML listed both
hooks in an isolated normal checkout; separate user-config probes also accepted
both TOML and JSON. Changing formats is therefore not an established fix.
Merge the hook configuration and update the primary checkout, then start a new
CLI session and check `/hooks` to trust it. Recheck configuration resolution
after upgrading Codex; this is an observed version-specific behavior, not a
guarantee for every release.

The shell hook launches a transient user systemd service, which runs
`bazel --noblock_for_lock clean --config=agent --expunge_async` in each built,
tracked Bazel workspace in the linked Git worktree. It skips the primary
checkout, submodules, workspaces without a `bazel-out` symlink, and busy Bazel
instances. Source files, Git worktrees, and shared disk caches are preserved.
Nested workspaces are discovered from tracked `MODULE.bazel`, `WORKSPACE`, and
`WORKSPACE.bazel` files. Custom convenience-symlink prefixes are not supported.

This Linux hook requires Bash, Git, Bazelisk's `bazel` on PATH, and a running
user systemd manager. The service inherits PATH, not the session's full
environment. `KillMode=process` allows Bazel's asynchronous deletion child to
finish after the service exits. Failures appear in the user journal; failure
to enqueue is reported by Codex as a hook failure. There is no automatic retry
or orphan-cache sweep. Concurrent sessions sharing a worktree may lose warm
build state when one ends; the Bazel lock protects running builds, not idle
sessions. Resuming a cleaned session rebuilds its outputs normally.

The hook uses shell because it must run from a fresh checkout without first
building a helper. It invokes Bazel directly to supply the startup lock flag,
which `bazel_agent` does not currently accept. Tests replace Git, systemd, and
Bazel with fixtures and never clean real caches.

Codex's `SessionEnd` hook is synchronous even with `async: true` and permits
at most three seconds. Systemd owns the longer-running cleanup; the hook does
not wait for it. See the [Codex hook documentation](https://learn.chatgpt.com/docs/hooks).
