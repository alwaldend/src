## Ownership

This design supersedes the rental-management direction in the scaffold's
[design](../add-mlflow-vastai-backend/design.md). Operators use Terraform to
provision and release rentals. Reuse the image-generation provisioning
components where they fit and configure a worker image with SSH and the
training runtime. The job server discovers existing rentals through the
repository's Vast.ai client.

The server owns scheduling and assignment. Each worker owns the execution
result for its assigned jobs. The MLflow Tracking Server owns experiment
metadata and artifacts.

## Interfaces

The Python entry point registers the backend as `vastai`. Its adapter uses
MLflow's project utilities to resolve the project, prepare its entry point,
and create a tracking run. It submits the prepared workload by gRPC and
exposes the tracking run ID through `SubmittedRun`. The workload runs with
the same MLflow run ID, experiment ID, and a tracking URI reachable from the
rental. It logs its own metrics and artifacts directly to MLflow.

The Go binary provides a job server and a worker service. The Python adapter
calls the job server's `SubmitJob`, `GetJob`, and `CancelJob` RPCs. The job
server invokes the binary's worker client over SSH. The worker client forwards
the request to the resident worker service by gRPC over a Unix socket. The
SSH invocation carries the request and response over stdin and stdout, with
diagnostics on stderr. The client and service use the same protobuf messages
for worker operations. The training command's stdout and stderr go to per-job
files.

## State and persistence

The adapter reuses one job ID across retries for the same MLflow submission.
The server persists the job ID and MLflow run ID before dispatch. It persists
the selected Vast.ai instance ID before requesting execution. The worker uses
the same job ID as its submission key and persists its local reservation before
launch. A retry of the same submission returns the existing job. Reusing an
ID with a different workload is an error.

The server and each worker use separate SQLite databases. The server stores
request and assignment data and the last confirmed worker state. The worker
records job state, process identity, log locations, and exit results. Each
service owns writes to its own database. The worker service outlives SSH
clients. Reconciliation after interruption uses the recorded assignment and
queries that same worker. If it is unreachable, the server retains the last
confirmed state and reports the communication failure separately. If a worker
or its storage is gone and no terminal result was observed, the outcome is
unknown.

Use the pinned `github.com/qmuntal/stateless` library for all server and worker
business decisions. Model request admission, rental eligibility and assignment,
submission, launch, observation, reconciliation, cancellation, and terminal
results as triggers, guards, transitions, and actions. Represent durable state
in SQLite, reconstruct machines from it, and use transactions to serialize
competing updates. The machine authorizes external effects after durable intent
is recorded. Reconcile ambiguous effects before retrying them. Input decoding,
gRPC, SSH, filesystem, database, and MLflow adapters implement I/O and protocol
translation. Inject `afero.Fs` for Go-owned file operations, following the
Anki as Code pattern.

The worker commits the terminal result before reporting it. A confirmed exit
code of zero maps to `FINISHED`, and a launch error or nonzero exit maps to
`FAILED` unless the worker has confirmed cancellation. Use `KILLED` for a job
cancelled before launch or confirmed stopped after cancellation. MLflow's
synchronous call waits on the server; disconnecting that caller leaves the job
under worker control. The adapter surfaces an unavailable or unknown outcome
as an error rather than inventing a terminal MLflow status.

The server reconciles jobs even without a waiting MLflow client. It projects
confirmed terminal results to the tracking run when the run is still active.
It retries a failed tracking update from its durable state without repeating
workload execution.
