# Autoscroll Specification

## Purpose

Turn pointer movement into configurable horizontal and vertical scrolling, with
configuration reloads during execution. This baseline records checked-in
behavior at revision `550d7e79b1f5fdbc2b6017b75178471d6914082f`, observed on
2026-09-08; desktop integration was not exercised.

Sources: [project README](../../../README.md),
[scroll controller](../../../main/py/autoscroll/_internal/autoscroll.py),
[scrolling and configuration support](../../../main/py/autoscroll/_internal/support.py),
and [argument parser](../../../main/py/autoscroll/_internal/arguments.py).

## Requirements

### Requirement: Button-controlled scrolling

Autoscroll SHALL start scrolling from the pointer position recorded when the
configured start button is pressed and stop on the configured end button or
release of the start button when hold mode is enabled.

#### Scenario: Release the held start button

- **WHEN** scrolling is active with hold mode enabled and the start button is
  released
- **THEN** Autoscroll SHALL stop scrolling.

### Requirement: Movement-dependent direction and interval

Autoscroll SHALL emit zero movement inside the configured square dead area and
derive scrolling direction from the pointer's displacement from its starting
position. Pointer movement SHALL recalculate the interval using the configured
speed, acceleration, and maximum absolute displacement.

#### Scenario: Return to the dead area

- **WHEN** both absolute coordinate displacements are no greater than the
  configured dead-area distance
- **THEN** the next scrolling direction SHALL be zero on both axes.

### Requirement: Configuration reload

When configuration monitoring is enabled, Autoscroll SHALL poll the configured
file at the configured interval and parse updated command-line arguments when
the file's modification time changes.

#### Scenario: Update runtime arguments

- **WHEN** a monitored configuration file changes to contain valid arguments
- **THEN** the next polling cycle SHALL apply those arguments to the running
  configuration without requiring a process restart.
