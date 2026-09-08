# Android Launcher Specification

## Purpose

Provide a text-based Android home application with app visibility controls and
local settings. This baseline records checked-in behavior at revision
`550d7e79b1f5fdbc2b6017b75178471d6914082f`, observed on 2026-09-08; it does not
record an on-device verification.

Sources: [project README](../../../README.md),
[manifest](../../../AndroidManifest.xml),
[home rows](../../../main/java/ui/HomeRow.kt),
[home actions](../../../main/java/ui/HomeViewModel.kt),
[launcher manager](../../../main/java/LauncherManager.kt),
and [state repository](../../../main/java/LauncherStateRepository.kt).

## Requirements

### Requirement: Text application launching

The launcher SHALL render application labels as clickable text and invoke the
selected package's Android launch intent.

#### Scenario: Select an application

- **WHEN** a user taps an application row with an available launch intent
- **THEN** the launcher SHALL start the selected package through the Android
  package manager.

### Requirement: Persistent application visibility

The launcher SHALL store hidden status per package, omit hidden applications
unless the show-hidden setting is enabled, and preserve hidden status when
refreshing the installed application list.

#### Scenario: Refresh a hidden application

- **WHEN** the installed application list is reloaded and a previously hidden
  package remains installed
- **THEN** its hidden status SHALL remain set and its row SHALL remain omitted
  while show-hidden is disabled.

### Requirement: Android home integration without network permission

The application manifest SHALL register the main activity for the Android HOME,
DEFAULT, and LAUNCHER categories without requesting the INTERNET permission.

#### Scenario: Inspect launcher registration

- **WHEN** the application manifest is inspected
- **THEN** its main activity SHALL be eligible for home application selection
  and the manifest SHALL contain no INTERNET permission request.
