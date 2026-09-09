---
title: ActivityWatch ingester Android
description: Accessibility-based Android collector that sends browser URLs to ActivityWatch
statuses:
  - maintenance
languages:
  - kt
tags:
  - android
  - activitywatch
---

# ActivityWatch ingester Android

An Android collector that reads supported browser address-bar text through
Android accessibility and sends it to an ActivityWatch server. The app holds
full accessibility permission only to obtain the URL; it never uses that
permission for anything else.

## How it works

- Enabling accessibility for this app lets it read the focused address bar.
- The accessibility service matches browser view IDs (`url_bar` and common
  Chrome-derived variants), extracts the current URL, and sends an
  ActivityWatch heartbeat for the `web.tab.current` bucket to a configurable
  server.
- `web.tab.current` heartbeats use `pulsetime=1.0`; consecutive URLs are sent
  with the time spent on the previous tab as the `duration` field.
- Data goes directly from this app to the ActivityWatch HTTP API; it is not
  uploaded to any other service.

## Configuration

Open the app to configure:

- **ActivityWatch server URL** — the base URL of the ActivityWatch server.
  The default is `http://127.0.0.1:5600` for the ActivityWatch Android app on
  the same device. Changes are saved immediately.
- **API key (optional)** — bearer token sent as `Authorization: Bearer ...`
  when the ActivityWatch server has API authentication enabled.

The `INTERNET` permission is required to reach the ActivityWatch HTTP API. The
default loopback URL is sent over cleartext HTTP only to `127.0.0.1` and
`localhost`, so no host data leaves the device.

## Supported browsers

Chrome and Chromium-based browsers that expose the `url_bar` view ID
(including ordinary Chrome and many Chrome-based forks). Firefox, Samsung
Internet, Opera, and Edge do not use that view ID yet.

## Build and install

1. Build the APK:
   `bazel build //projects/activitywatch_ingester_android/main/java:ingester_binary`
2. Verify the APK starts on an already running device or emulator:
   `bazel run //projects/activitywatch_ingester_android/test:ingester_smoke_test`
   This manual smoke test installs the Bazel-built APK, launches the activity,
   waits for the process, and confirms it stays alive. It is not a CI test.
3. Install the APK on the device from `bazel-bin/projects/activitywatch_ingester_android/main/java/ingester_binary.apk`.
4. Start the app and set the ActivityWatch server URL.
5. Enable accessibility for the app in Android settings; Android shows the
   full-access warning, which is the mechanism this collector is built on.
6. Open a browser. The current URL is re-sent to ActivityWatch whenever the
   address bar changes.

## Compatibility

- Minimum Android SDK 25 (Android 7.0).
- Tested against the ActivityWatch v0.13 HTTP API (`POST /api/0/buckets` and
  `POST /api/0/buckets/web.tab.current/heartbeat`).

## Landing page

The landing page is published at
[activitywatch-ingester-android.alwaldend.com](https://activitywatch-ingester-android.alwaldend.com/),
built from this README by the shared Hugo landing template. The DNS record
lives in `dnsconfig.json` and is applied through `//infra/dns`.
