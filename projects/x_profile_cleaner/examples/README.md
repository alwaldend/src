# Archive script examples

The generic scripts below use the worker's `Evaluation.script` entrypoint. Each
call adds an operation to the plan; evaluation performs no X requests. The current
backend executes post deletion and unrepost. The other examples demonstrate typed
planning for additional executor implementations.

`policy.lua` uses the cleanup DSL accepted by the CLI/browser. Its test case selects
the cleanup entrypoint, which wraps it through the same application helper.

| Example                                                | What it demonstrates                                                                                                                                                   |
| ------------------------------------------------------ | ---------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| [operations.lua](operations.lua)                       | Remove every exported like and following relationship.                                                                                                                 |
| [delete_old_posts.lua](delete_old_posts.lua)           | Select posts/replies older than 90 days with fewer than ten known likes; retain configured IDs, media, unknown counters and counters beyond Lua's exact numeric range. |
| [remove_old_reposts.lua](remove_old_reposts.lua)       | Select repost relationships older than 180 days, skip unresolved sources, and deduplicate original post IDs.                                                           |
| [prune_relationships.lua](prune_relationships.lua)     | Match liked-post text without Lua patterns and retain an explicit following list. Missing text does not match.                                                         |
| [read_referenced_posts.lua](read_referenced_posts.lua) | Plan explicit reads for unique reply/quote targets using typed reference kinds.                                                                                        |
| [profile_announcement.lua](profile_announcement.lua)   | Read account metadata and activity count, handle a missing username, and order a profile update before post creation.                                                  |
| [policy.lua](policy.lua)                               | Use the cleanup DSL for age selection, preservation and explicit forced deletion.                                                                                      |

Edit thresholds, keywords and retained IDs before submitting a script. Generic
deletion selects individual activities; use the cleanup DSL when reply/thread
preservation is required. Only files and metadata present in the export are
available, so an archive is not a current view of the account.

## Automated acceptance

Run from the repository root:

```sh
bazel test //projects/x_profile_cleaner/test/worker:examples_test
```

The [JSON configuration](../test/worker/examples/cases.json) owns each case's
archive input, fixed evaluation time, example, entrypoint and complete expected
protobuf plan. Raw YTD fixture data is stored as text and mapped to standard
`data/*.js` ZIP entries by that configuration. The test launches the real worker,
submits each script over gRPC and compares its result using protobuf equality,
including operation order, optional fields and cleanup reports.

Cases cover populated/empty exports, missing metadata, exact large IDs, media and
age exclusions, retained IDs and duplicate targets. Every Lua file in this
directory must have a configured case. Example results are preserved under
`bazel-testlogs/projects/x_profile_cleaner/test/worker/examples_test/test.outputs/`.
The target is also included in `//projects/x_profile_cleaner:test`.
