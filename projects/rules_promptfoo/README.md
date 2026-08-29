---
title: Rules Promptfoo
description: A pinned Bazel runner for Promptfoo skill evaluations
languages:
  - bzl
  - javascript
tags:
  - bzl_rules
  - llm_evals
---

`rules_promptfoo` runs repository skills through a pinned Promptfoo CLI. It
keeps dependency resolution reproducible in Bazel while deliberately running
model calls only at test runtime.

## Setup

Add the module and its skill-provider dependency:

```starlark
bazel_dep(name = "rules_promptfoo", version = "<VERSION>")
bazel_dep(name = "rules_skill", version = "<VERSION>")
```

The development checkout uses sibling `local_path_override` declarations.
Published releases must replace those local overrides with registry versions.

## Evaluate a skill

```starlark
load("@rules_promptfoo//promptfoo:defs.bzl", "promptfoo_test")

promptfoo_test(
    name = "answer_question_eval",
    config = "promptfooconfig.yaml",
    skills = ["//projects/agents/skills/answer-question:skill"],
    env_inherit = [
        "CODEX_HOME",
        "CODEX_PATH_OVERRIDE",
    ],
    reuse_codex_login = True,
)
```

The runner physically copies each selected `SkillInfo.files_by_path` bundle
to an isolated Bazel test-temporary
`workspace/.agents/skills/<name>/` tree. Point the Codex SDK provider at the
exported workspace and, when using a compatible host Codex executable, at the
exported executable path:

```yaml
prompts:
  - "{{ question }}"
providers:
  - id: openai:codex-sdk
    config:
      codex_path_override: "{{ env.CODEX_PATH_OVERRIDE }}"
      working_dir: "{{ env.PROMPTFOO_SKILL_WORKSPACE }}"
      cli_env:
        CODEX_HOME: "{{ env.PROMPTFOO_SUBJECT_CODEX_HOME }}"
tests:
  - vars:
      question: Why is the sky blue?
    assert:
      - type: skill-used
        value: answer-question
```

Live eval targets are tagged `manual`, `requires-network`, `no-cache`,
`no-remote`, and `local`. They also pass `--no-cache`, `--no-write`, and
`--no-share` to Promptfoo; `local` prevents sandboxed or remote execution when
the caller opts into host Codex credentials. The runner exports a separate
empty `PROMPTFOO_JUDGE_WORKSPACE` so a model judge need not discover the skill
it grades. Explicit results are preserved as the Bazel undeclared output
`results.json`; Promptfoo configuration, cache, isolated workspaces, and
subprocess temporary files stay in a mode-0700 directory below Bazel's
absolute `TEST_TMPDIR` and are removed on exit. Bazel owns and cleans that
private parent directory as a second cleanup boundary. Callers that must
redirect test scratch can use Bazel's `--test_tmpdir` option.

The runner deliberately uses Bazel-managed temporary storage instead of a
source-checkout `out/` directory. Putting a no-skill control below the source
tree would expose the checkout's `.agents` directory through an ancestor and
invalidate the isolation that the control is meant to measure. This is the
repository policy's operating-system-temporary exception for a tool that
cannot safely use task-local source storage; the runner still removes its
randomized child directory on normal exit and handled signals. The runner
rejects a `TEST_TMPDIR` below any ancestor containing `.agents`, so do not
redirect `--test_tmpdir` into a source checkout or another agent tree.

`results.json` contains the evaluated prompts and model outputs. CI may collect
undeclared outputs when a manual eval is explicitly run, so treat the artifact
as potentially sensitive even though the target is `no-cache` and `no-remote`.

Run a live target explicitly and pass the absolute path to an existing Codex
login and a compatible host Codex executable through the test environment,
never credential contents through `env` or a checked-in config. Set the
provider's `codex_path_override` from `CODEX_PATH_OVERRIDE`:

```sh
bazel_agent test //path/to:answer_question_eval \
  --test_env=CODEX_HOME=/absolute/path/to/.codex \
  --test_env=CODEX_PATH_OVERRIDE=/absolute/path/to/codex \
  --test_output=errors
```

If that Codex installation reaches OpenAI through host proxy variables, pass
only the needed names as well (for example, `--test_env=HTTPS_PROXY` and
`--test_env=NO_PROXY`) and declare the same names in the rule's `env_inherit`.
Do not put proxy values in checked-in files.

With `reuse_codex_login = True`, the runner links only the writable
`CODEX_HOME/auth.json` into distinct, otherwise empty subject and judge Codex
homes inside its mode-0700 state directory. Provider configs can select them
through `PROMPTFOO_SUBJECT_CODEX_HOME` and `PROMPTFOO_JUDGE_CODEX_HOME`. The
runner holds an exclusive advisory lock on the stable login directory for the
complete test, serializing repository eval targets that use the same login.
With a compatible Codex executable, an automatic refresh writes through either
link to the persistent host file, while the host's `config.toml`, skills,
plugins, memories, and MCP configuration remain unstaged and cannot contaminate
the subject, no-skill control, or judge. This reuses the caller's existing
login state; it does not mint or export a token. An API-key-backed target may
instead inherit `OPENAI_API_KEY` without enabling login reuse.

For this purpose, compatible means that file-backed auth updates open and
truncate `$CODEX_HOME/auth.json`, preserving the symlink. The bundled Codex
0.144.0 implementation and the tested host Codex 0.150.1 do this. Do not
override the executable with an implementation that atomically replaces the
isolated `auth.json` path: that would replace the symlink, split subject and
judge state, and fail to persist a rotated refresh token. See the corresponding
OpenAI Codex file-storage implementations for
[0.144.0](https://github.com/openai/codex/blob/rust-v0.144.0/codex-rs/login/src/auth/storage.rs)
and
[0.150.1](https://github.com/openai/codex/blob/rust-v0.150.1/codex-rs/login/src/auth/storage.rs).

Login reuse also forces Promptfoo's CLI, assertion-worker, and prompt-suggestion
concurrency to one, overriding higher configuration values so subject and judge
Codex processes cannot race the same single-use refresh token within one test.

Do not run another Codex process against the same `auth.json` concurrently;
unrelated processes do not honor the runner's advisory lock. ChatGPT-managed
auth can rotate its refresh token, so the refreshed file must remain persistent
and have a single serialized user. The local-login mode is unsuitable for
shared or untrusted CI; use an API key or supported workload identity there.
In this public repository it is strictly a manual, trusted-developer workflow;
never seed or persist ChatGPT-managed auth in repository CI or artifacts.
See OpenAI's
[managed-auth guidance](https://learn.chatgpt.com/docs/auth/ci-cd-auth) for the
refresh and serialization requirements.

Run credentialed live targets only from a reviewed, trusted revision. The
isolated working directory and read-only provider setting reduce accidental
writes, but Codex still has whatever host read access its enforced sandbox
policy permits and can send prompt context to the configured service.

Files in `data` are placed in the test runfiles. `args` are passed literally;
they do not expand Bazel location expressions. Prefer keeping checked-in
Promptfoo prompt and dataset files beside the config and list them in `data`
so Promptfoo's own config-relative path resolution remains portable.
The `env` and `env_inherit` attributes cannot override Bazel's test temporary
directories or the runner-owned `TMPDIR`, `TMP`, and `TEMP` values.

## Validate without a model call

```starlark
load(
    "@rules_promptfoo//promptfoo:defs.bzl",
    "promptfoo_validate_test",
)

promptfoo_validate_test(
    name = "promptfoo_config_test",
    config = "promptfooconfig.yaml",
)
```

Validation is an ordinary offline Bazel test and does not receive the live
eval tags. It rejects `env_inherit` so validation cannot accidentally receive
host credentials.

## Dependency boundary

The module pins Promptfoo and the Codex SDK with a dedicated pnpm lock. npm
lifecycle hooks and optional dependencies are disabled. The Codex executable
and Promptfoo's libSQL binding are promoted to explicit Linux x86-64
dependencies. AJV is also promoted because `ajv-formats` declares its runtime
peer as optional, which would otherwise be removed by `no_optional`. The
bundled real CLI target is intentionally compatible only with glibc-based
Linux x86-64. Bazel enforces the OS and CPU constraints, but this repository
does not currently expose a libc constraint. Other Promptfoo providers that
rely on omitted optional packages or install scripts must be reviewed and
promoted explicitly before use.

The pinned Promptfoo distribution is patched so
`PROMPTFOO_DISABLE_TELEMETRY=1` also suppresses its best-effort request to
`r.promptfoo.app`. The upstream 0.122.2 telemetry implementation otherwise
records a "telemetry disabled" event through that separate endpoint. The same
patch makes Promptfoo's otherwise hard-coded prompt-suggestion concurrency
honor the runner's serialization setting.
