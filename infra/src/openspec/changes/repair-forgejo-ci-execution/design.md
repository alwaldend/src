## Context

Run 38 started persistent Bazel servers across workspaces. At 08:19:13 UTC the
runner host exhausted memory and killed a Java process in the Forgejo service
that was using about 4 GiB of resident memory. The job received SIGTERM within
the same second. The two earlier startup failures did not include a JVM cause.

## Decisions

Proceed with shutting down the server after each workspace's test attempt on
Forgejo. This limits the lifetime of the CI-owned servers while reusing them
between build and test. The stable checkout and disk caches continue to serve
later jobs. Each workspace starts a new server on the next job.

Keeping all servers resident preserves warm analysis and leaves the confirmed
memory pressure unaddressed. Batch mode starts a JVM for each phase. Increasing
host RAM would require an infrastructure change and would leave idle servers
resident as the workspace count grows.

## Risks and Trade-offs

An individual workspace can still exceed host memory. Verify aggregate memory
behavior on the next trusted Forgejo run.
