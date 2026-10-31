## Decisions

Reject unrecognized fields at each owned structure rather than accepting obsolete rules as a compatibility alias. Silent fallback is the defect; producer and consumers are updated together. Downloader and execute keyword dictionaries remain governed by Bazel's repository APIs.

## Failure cases and validation

Before implementation, add executable analysis fixtures for obsolete fields, unknown top-level and nested fields, malformed constraint metadata, and preservation of Linux/x86-64 constraints. Existing runtime-files fixtures must continue to pass. Validate repository locks and representative generated toolchain declarations without changing their selected archives.

## Continuation

The shared integration change owns delivery. Evidence belongs under out/repository-organization. Next action: write analysis fixtures before parser changes.
