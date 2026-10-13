## Purpose

Provide reproducible taxi-fare training with pinned data and verifiable saved artifacts.

## ADDED Requirements

### Requirement: Dependency-backed training

The lesson SHALL train using the existing pinned upstream CSV and repository-locked
Python dependencies without a checked-in CSV copy.

#### Scenario: Build training artifacts

- **WHEN** the user builds the lesson model target
- **THEN** sandboxed training produces model.pkl and metrics.json from the full dependency CSV
- **AND** metrics identify the split, model parameters, dataset, and finite R²

### Requirement: Reproducible executable

The training executable SHALL accept dataset, output directory, split fraction,
and seed parameters, exclude identifier and target columns from features, and
filter missing required values and zero coordinate pairs.

#### Scenario: Repeat a bounded training run

- **WHEN** the same upstream prefix is trained twice with the same seed and split
- **THEN** saved estimators produce equal predictions and equal finite R² values
- **AND** the end-to-end test retains a verification artifact
