# Bazel scratch-isolation evidence

## Subject

Base commit `6780d53a69e32064d648e6a04f1c0cecd7d713fd`, the accepted
`phase1-contracts-001` source set, and this attempt's exact files:

| SHA-256 | Path |
| --- | --- |
| `3661c61662874e0948976c4435bc55e93d20f5bcb3c3f35bcf816d262c3273fb` | `projects/bazel_agent/README.md` |
| `dcdef462cb610b277bd117c212fb74d004d596a301b6142fb63b7bba2786c614` | `projects/bazel_agent/cmd/bazel_agent/main.go` |
| `731717430b1e2854a4facbea76b6ad2552a0bf0fb273517a51d604857f42c1ee` | `projects/bazel_agent/cmd/bazel_agent/main_test.go` |
| `2a80052ce5bf62d40feb77a072ec81e54e6158a110ed0b8216522a3231b4a429` | `tools/bazelrc/BUILD.bazel` |
| `274e275ef6b85a8f14bccaef3ca19c1d5e1629bc776afe5e69f130f71bcb78aa` | `tools/bazelrc/README.md` |
| `12026e5f459ef4c951ce17471dd856202e3ef805a4a6a0fb949152f941c815f8` | `tools/bazelrc/agent_environment_test.sh` |
| `aa276078d84775823bd9ae757d3ac08c3ab3651ea214dc75458eaf6d8eee9bd6` | `tools/bazelrc/project.bazelrc` |

## Verification

- `bazel_agent test //projects/bazel_agent/cmd/bazel_agent:all
  //tools/bazelrc:all //tools/agents/api/v1alpha1:all //projects/goal/...
  //:buildifier_test` passed all 11 tests.
- `bazel_agent build --config=lint //projects/bazel_agent/...
  //tools/bazelrc/... //tools/agents/api/v1alpha1:all
  //projects/goal/internal/fsstore:all` built all 17 selected targets and lint
  aspects successfully.
- `git diff --check` passed.
- Runner unit coverage proves the exact ambient environment is preserved.
- `//tools/bazelrc:agent_environment_test` rejects any future `agent` profile
  propagation of `TEMP`, `TMP`, or `TMPDIR` into repository rules, actions,
  host actions, or tests.

## Criterion observations

- `scratch-isolation`: fail overall, with material improvement. The mandatory
  Bazel path no longer creates shared `out/tmp` or injects an absolute host
  path into execution environments. Worktree-global Cordis scratch and the
  missing task/run manifest remain unresolved.
- `exact-candidate-validation`: pass for this attempt's affected source and
  all fixed regressions.
- `shared-contracts`, `information-policy`, and `encountered-bug-policy` remain
  passing against the cumulative candidate because their full prior test set
  was rerun.

The host-installed `bazel_agent` binary was not replaced because this task did
not authorize changing the agent environment. The built target contains the
new behavior; repository delivery must not claim host installation.
