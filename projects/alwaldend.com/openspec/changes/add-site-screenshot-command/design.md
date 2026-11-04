## Context

The site already owns Puppeteer and pinned Chrome Headless Shell for browser
tests. Shared preview automation timed out during the requested captures. A
scratch capture proved that request interception can serve built files directly
to the browser without an HTTP server.

## Goals / Non-Goals

Make the repeated screenshot operation a discoverable Bazel command. Capture
existing presentation; do not change page content, hosting, or sorting.

## Decisions

Use a small Node command beside the existing search command, with Bazel supplying
the site and browser paths. Reuse the established Puppeteer import and runfiles
model. Standard argument parsing and local asset interception need no dependency.
The command resolves requests beneath the site directory, blocks external
origins, waits for fonts, and records page titles and PNG dimensions. It closes
the browser and removes the private temporary profile on completion or failure.

## Failure cases before implementation

Reject invalid dimensions, color schemes, nonlocal page paths, filename
collisions, and missing pages. Prevent asset requests from escaping the site
directory. Browser failures must produce a nonzero exit and clean up the profile.
An offline command fixture verifies real PNG output, full-page height, CSS asset
loading, custom viewport/theme options, external-request blocking, and failures.

## Risks / Trade-offs

Browser rendering is platform-dependent, so the manifest records configuration
and dimensions rather than promising identical PNG bytes across platforms.
The command requires no running server, browser login profile, or deployment.
