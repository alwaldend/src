## 1. Integrate the dependency

- [x] 1.1 Add the pinned executable lock, module fragment, alias and README;
      verify the SRI against the publisher checksum and query the package targets.
- [x] 1.2 Regenerate the root module and reconcile Bzlmod dependencies;
      inspect the generated diff for unrelated changes.

## 2. Validate and deliver

- [x] 2.1 Build the package and run its alias with `--version` and
      `--list-extractors`; preserve the pinned version and extractor output under
      ignored `out/add-yt-dlp/`.
- [x] 2.2 Complete formatting, applicable repository quality checks and OpenSpec
      validation against the exact candidate; archive the completed change.
