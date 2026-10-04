## Purpose

Make the two upstream course datasets available to the MLOps application at
runtime with reproducible source identity and no checkout-path dependency.

## ADDED Requirements

### Requirement: Package pinned upstream course data

The application runtime bundle SHALL include data.csv and dataset.csv from
an immutable, integrity-verified revision of the user-requested
Yandex-Practicum/mlops-freetrack repository. Both files SHALL be available
through application runfiles without network downloads or dependence on the
repository working directory. Their bytes SHALL match the pinned upstream files.

#### Scenario: Launch outside the repository working directory

- **WHEN** the Bazel-built application is launched from a different directory
- **THEN** both CSV files resolve through its runfiles and can be opened
- **AND** their checksums match the pinned upstream source

#### Scenario: Check health with packaged course data

- **WHEN** the server starts with both course CSV files in its runtime bundle
- **THEN** GET /health retains its agreed response without parsing the datasets
  or invoking recommendation computation

### Requirement: Build the course model artifact

The repository SHALL expose a runnable upstream training target and a Bazel
build target whose declared output is model.pkl. Training SHALL use the pinned
course CSV files and resolved upstream Python dependencies without network
access or writing into source directories. The artifact SHALL be usable by a
future model consumer; the recommendation server SHALL load it through the upstream loader.

#### Scenario: Build the trained model

- **WHEN** the model build target is selected
- **THEN** it runs the pinned upstream training code against the pinned datasets
  and produces a nonempty model.pkl in the Bazel output tree
- **AND** a verification consumer loads it and obtains a structurally valid
  recommendation from a track present in the model
