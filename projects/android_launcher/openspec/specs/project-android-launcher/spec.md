# Android Launcher Specification

## Purpose

Provide a text-based Android home application that launches installed
packages, hides selected applications, and keeps its presentation in local
settings. It runs without ads, tracking, or network access. This baseline
records checked-in behavior; on-device verification is recorded per change.

Sources: [project README](../../../README.md),
[manifest](../../../AndroidManifest.xml),
[contracts](../../../main/proto/contracts.proto),
[home route](../../../main/java/ui/HomeRoute.kt),
[home rows](../../../main/java/ui/HomeRow.kt),
[home dialog](../../../main/java/ui/HomeDialog.kt),
[home view model](../../../main/java/ui/HomeViewModel.kt),
[launcher manager](../../../main/java/LauncherManager.kt),
[state repository](../../../main/java/LauncherStateRepository.kt),
[settings screen](../../../main/java/ui/SettingsRoute.kt),
and [theme](../../../main/java/ui/LauncherTheme.kt).

## Requirements

### Requirement: Text application launching

The launcher SHALL render application labels as clickable text and invoke the
selected package's Android launch intent.

#### Scenario: Select an application

- **WHEN** a user taps an application row with an available launch intent
- **THEN** the launcher SHALL start the selected package through the Android
  package manager.

### Requirement: Installed application discovery

The launcher SHALL enumerate installed launcher activities through the package
manager, storing each package's label, and SHALL register a `LauncherApps`
callback so package additions, removals, and shortcut changes refresh its
stored application list.

#### Scenario: Install or remove a package

- **WHEN** Android reports an added, removed, or changed package
- **THEN** the launcher SHALL add, drop, or refresh that package's row without
  a restart.

### Requirement: Persistent application visibility

The launcher SHALL store hidden status per package, omit hidden applications
unless the show-hidden setting is enabled, and preserve hidden status when
refreshing the installed application list.

#### Scenario: Refresh a hidden application

- **WHEN** the installed application list is reloaded and a previously hidden
  package remains installed
- **THEN** its hidden status SHALL remain set and its row SHALL remain omitted
  while show-hidden is disabled.

#### Scenario: Toggle hidden applications from the launcher

- **WHEN** the user chooses to show or hide hidden applications from the home
  dialog
- **THEN** the launcher SHALL persist the new visibility preference and update
  the rows immediately.

### Requirement: Application long-press actions

Long-pressing an application label SHALL open a dialog showing the application
label, its icon when the stored icon decodes, its launcher shortcuts when the
launcher owns the home role, and actions to open system app info, toggle
hidden status, or uninstall the package. The pressed row SHALL show visual
feedback while the gesture is held, and the row whose dialog is open SHALL
stay highlighted.

#### Scenario: Long-press an application

- **WHEN** the user long-presses an application label
- **THEN** the launcher SHALL highlight that row and open its action dialog.

#### Scenario: Icon arrives after the dialog opens

- **WHEN** the stored application entry gains an icon while its dialog is open
- **THEN** the dialog SHALL render the decoded icon without restarting the
  activity or crashing on the empty placeholder.

### Requirement: Launcher shortcuts and dialogs

The launcher SHALL read dynamic, pinned, and manifest shortcuts for a package
when the app dialog loads and SHALL present them as dialog buttons that start
the selected shortcut. The dialog SHALL omit the shortcut section outside the
home role and SHALL dismiss after an action.

#### Scenario: Launch an application shortcut

- **WHEN** the app dialog lists a shortcut and the user selects it
- **THEN** the launcher SHALL start that shortcut through `LauncherApps`.

### Requirement: Local settings and layout

The launcher SHALL persist its settings in a protobuf DataStore, including
application-card label transformation, text style, text color, font family and
padding, plus layout arrangement, sort order, and screen padding. Settings
changes SHALL apply to the home rows without restarting.

#### Scenario: Change a layout setting

- **WHEN** the user selects a different arrangement, padding, or sort order
- **THEN** the home rows SHALL re-render with the stored setting.

### Requirement: Settings backup and reset

The settings screen SHALL export the stored state to a user-selected document,
import a previously exported state stream, and reset settings to their
defaults without discarding application state.

#### Scenario: Reset settings

- **WHEN** the user confirms the reset action
- **THEN** the launcher SHALL restore default card and layout settings while
  keeping installed application entries.

### Requirement: Themed presentation

The launcher SHALL render its surface through a Material 3 theme that follows
the system light and dark appearance and uses platform dynamic colors on
Android 12 and newer, and SHALL keep content inside the system window insets.

#### Scenario: Inspect launcher theme

- **WHEN** the launcher renders on a device in dark appearance
- **THEN** its surface, rows, and dialogs SHALL use the dark color scheme.

### Requirement: Android home integration without network permission

The application manifest SHALL register the main activity for the Android HOME,
DEFAULT, and LAUNCHER categories, SHALL declare the delete-package permission
used by the uninstall action, and SHALL NOT request the INTERNET permission.

#### Scenario: Inspect launcher registration

- **WHEN** the application manifest is inspected
- **THEN** its main activity SHALL be eligible for home application selection
  and the manifest SHALL contain no INTERNET permission request.

#### Scenario: Report the missing home role

- **WHEN** the launcher is not the selected home application
- **THEN** the home screen SHALL show a banner explaining that shortcuts are
  unavailable.
