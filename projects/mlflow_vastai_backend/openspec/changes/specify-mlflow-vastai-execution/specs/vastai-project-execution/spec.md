## ADDED Requirements

### Requirement: MLflow backend contract

The backend SHALL register as `vastai` under `mlflow.project_backend`, prepare
the requested project entry point with MLflow, create its tracking run before
submission returns, and expose the run ID through a `SubmittedRun`. The backend
SHALL pass the tracking URI, experiment ID, and run ID to the workload.

#### Scenario: Submit and wait for a project

- **WHEN** an MLflow caller submits a valid project to the Vast.ai backend
- **THEN** the backend submits the prepared job through gRPC and returns a handle with the MLflow run ID
- **AND** the handle can observe status, wait for success or failure, and request cancellation
- **AND** metrics and artifacts explicitly logged by the workload go to the associated MLflow run

### Requirement: Discover Terraform-managed workers

The job server SHALL use read-only Vast.ai discovery to resolve existing rentals
from configured instance IDs or exact worker labels. It SHALL assign only an
eligible, running rental with a usable SSH endpoint. Rental provisioning and
release SHALL remain under operator-controlled Terraform.

#### Scenario: No eligible worker is reachable

- **WHEN** discovery finds only unrelated, stopped, or otherwise ineligible rentals
- **THEN** the server reports that no eligible worker can accept the job
- **AND** it leaves the rentals unchanged

### Requirement: Communicate through bounded service interfaces

One Go binary SHALL provide a job server, a resident worker service, and a
short-lived worker client. The MLflow plugin SHALL call the job server by gRPC.
The job server SHALL invoke the worker client through SSH. The worker client
SHALL call the worker service by gRPC over a rental-local Unix socket.

#### Scenario: SSH caller exits after submission

- **WHEN** the worker accepts a job and the SSH client exits
- **THEN** the worker service continues to own the job
- **AND** the server can query or cancel it through a later SSH invocation

### Requirement: Drive all business logic with Stateless

All server and worker business logic SHALL be implemented through
`github.com/qmuntal/stateless` state machines. Domain events SHALL enter as
triggers. Business decisions and effects SHALL be expressed as guards,
transitions, or actions. Transport handlers, storage adapters, and the Python
plugin SHALL delegate domain decisions to these machines.

#### Scenario: Cancellation competes with job completion

- **WHEN** a cancellation request and a workload exit are observed for the same job
- **THEN** the worker processes them through the same persisted state machine
- **AND** a committed terminal result is preserved across subsequent requests and restarts

### Requirement: Persist and reconcile execution state

The job server and each worker SHALL use separate SQLite databases. The server
SHALL durably associate a job ID with its MLflow run and assigned rental. The
worker SHALL persist job identity, process state, and terminal results, and
retain separate stdout and stderr files. The worker SHALL be authoritative
for workload outcomes.

#### Scenario: Submission response is lost

- **WHEN** the worker accepts a job and the server loses the SSH response
- **THEN** the server retries or queries the same worker using the original job ID
- **AND** the worker returns the existing job without starting another workload

#### Scenario: The assigned rental becomes unreachable

- **WHEN** the server cannot reach the worker after dispatch
- **THEN** it retains the assignment and last confirmed state
- **AND** it reports the observation failure independently of the workload result
- **AND** it reconciles the job on the same worker when connectivity returns

#### Scenario: The rental is lost before the result is observed

- **WHEN** the worker and its stored execution result are no longer available
- **THEN** the server reports an unknown outcome if it has no confirmed terminal result
- **AND** it does not resubmit the job to another rental automatically

### Requirement: Report completion and cancellation accurately

The worker SHALL record a terminal result before making it available to the
server. The adapter SHALL map accepted and executing jobs to `SCHEDULED` and
`RUNNING`, successful jobs to `FINISHED`, confirmed failures to `FAILED`,
and confirmed cancellations to `KILLED`. The adapter SHALL report success from
`wait()` only for a confirmed successful result.

#### Scenario: Cancel an active training job

- **WHEN** MLflow requests cancellation of an active job
- **THEN** the worker stops the associated workload and records cancellation after confirming it stopped
- **AND** the rental remains available under Terraform management

#### Scenario: The submitting client exits before the job finishes

- **WHEN** an accepted job finishes after the Python client exits
- **THEN** the server independently reconciles the worker's durable terminal result
- **AND** it updates the MLflow tracking run if the run is still active
- **AND** a failed tracking update can be retried without relaunching the job
