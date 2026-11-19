# ActivityWatch Ingester Android Specification

## Purpose

Collect the focused browser address bar on an Android device through the
accessibility service and forward it to a local ActivityWatch server. The
application holds full accessibility access solely to read that address bar
text; this baseline records checked-in behavior, not an on-device verification.

Sources: [project README](../../../README.md),
[manifest](../../../AndroidManifest.xml),
[accessibility service](../../../main/java/ActivityWatchAccessibilityService.kt),
[HTTP client](../../../main/java/ActivityWatchIngester.kt),
[event model](../../../main/java/ActivityWatchEvent.kt),
[session tracker](../../../main/java/SessionTracker.kt),
[URL extractor](../../../main/java/UrlExtractor.kt),
[settings screen](../../../main/java/IngesterSettingsScreen.kt),
and [theme](../../../main/java/IngesterTheme.kt).

## Requirements

### Requirement: Accessibility address-bar capture

The application SHALL declare an accessibility service bound to
`android.permission.BIND_ACCESSIBILITY_SERVICE` that reports view IDs and
window content changes. The service SHALL resolve the focused address bar by
searching the event source for the browser's known URL bar view IDs and SHALL
reject empty results without sending a heartbeat.

#### Scenario: Recognized address bar changes

- **WHEN** an accessibility event arrives from a package whose source node
  contains a recognized URL bar view ID with non-blank text
- **THEN** the service SHALL treat that text as the current URL.

#### Scenario: Unknown or empty address bar

- **WHEN** the event source exposes no recognized URL bar view ID or only blank
  text
- **THEN** the service SHALL send no heartbeat for that event.

### Requirement: ActivityWatch heartbeat delivery

The application SHALL POST a `web.tab.current` bucket creation request for the
`aw-watcher-android-web` client of type `currenttab` on hostname `android`
before the first heartbeat, SHALL send heartbeats as
`POST /api/0/buckets/web.tab.current/heartbeat?pulsetime=1.0`, and SHALL strip
the `http://` or `https://` scheme from the reported URL. When the address bar
stays on one URL, the service SHALL include the elapsed time on that URL as the
heartbeat duration.

#### Scenario: Report a URL change

- **WHEN** the focused URL changes to a new value
- **THEN** the service SHALL send a heartbeat carrying the previous URL's
  elapsed duration before sending the new URL with zero duration.

#### Scenario: Keep the same URL active

- **WHEN** the same URL remains focused across successive events
- **THEN** the reported heartbeat SHALL carry the elapsed time since the URL
  became focused.

### Requirement: Bounded outgoing work

The ingester SHALL coalesce pending heartbeats to the latest URL so an
unreachable or slow server cannot accumulate an unbounded backlog, and SHALL
treat connection and read failures as non-fatal.

#### Scenario: Server is unreachable

- **WHEN** the configured ActivityWatch server refuses a connection or exceeds
  the client timeout
- **THEN** the failure SHALL remain local, the service SHALL stay enabled, and
  only the newest URL SHALL remain queued for delivery.

### Requirement: Configurable server connection

The application SHALL persist the ActivityWatch base URL and optional API key
in shared preferences and apply changes to the running service without a
restart. When an API key is configured, requests SHALL send it as an
`Authorization: Bearer` header. The application SHALL request the INTERNET
permission, which carries only loopback destinations, and SHALL permit
cleartext HTTP only for `127.0.0.1` and `localhost`.

#### Scenario: Change the server URL while the service runs

- **WHEN** the base URL or API key preference changes
- **THEN** subsequent heartbeats SHALL use the new value without restarting
  the service.

#### Scenario: Inspect cleartext network policy

- **WHEN** the network security configuration is inspected
- **THEN** cleartext HTTP SHALL be permitted for loopback hosts and denied for
  other domains.

### Requirement: Capture controls and settings surface

The main activity SHALL expose controls for enabling accessibility ingestion,
extracting URLs, and key logging, and SHALL persist each toggle. URL extraction
and key logging SHALL be disabled until accessibility ingestion is enabled, and
the screen SHALL offer a shortcut that opens the Android accessibility
settings.

#### Scenario: Enable a dependent capture control

- **WHEN** accessibility ingestion is disabled
- **THEN** the URL extraction and key logging controls SHALL remain disabled
  and SHALL not enable capture on their own.

#### Scenario: Inspect key logging scope

- **WHEN** key logging is enabled
- **THEN** the preference SHALL be stored and read by the service, and no key
  events SHALL be transmitted until a later specification defines them.

### Requirement: Themed system-bar presentation

The ingester SHALL render its settings inside a Material 3 surface that follows
the system light and dark appearance, using platform dynamic colors on Android
12 and newer. The status-bar region SHALL remain part of the themed
background, and the status-bar icon appearance SHALL contrast with the active
light or dark surface.

#### Scenario: Switch the system appearance

- **WHEN** the device changes between light and dark appearance
- **THEN** the ingester surface and status-bar icons SHALL update together
  without leaving an untinted system-bar band.
