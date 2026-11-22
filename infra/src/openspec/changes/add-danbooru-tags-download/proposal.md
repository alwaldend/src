## Why

Repository consumers need a reproducible download of the Danbooru tag vocabulary in CSV format.

## What Changes

Add a Bazel target for the full DeepGHS tag snapshot. Verify the fetched CSV against the upstream checksum and confirm that it has the expected columns.
