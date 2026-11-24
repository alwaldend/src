package jobrepository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"git.alwaldend.com/alwaldend/src/projects/mlflow_vastai_backend/api/jobpb"
	"google.golang.org/protobuf/proto"
)

// ErrNotFound indicates that the worker has no record for a job ID.
var ErrNotFound = errors.New("job not found")

// Record is the persisted representation of one worker job.
type Record struct {
	// ID is the MLflow run ID and SQLite primary key.
	ID string
	// Workload is immutable after insertion.
	Workload *jobpb.Workload
	// Phase holds the persisted job phase.
	Phase string
	// ExitCode holds the confirmed command exit code when available.
	ExitCode *int32
}

// Repository stores job records using a connection owned by the caller.
type Repository struct {
	// db provides the queries and transactions.
	db *sql.DB
}

// New borrows the caller's connection.
func New(db *sql.DB) *Repository {
	return &Repository{db: db}
}

// EnsureSchema creates the job table if it is missing.
func (repository *Repository) EnsureSchema(ctx context.Context) error {
	if _, err := repository.db.ExecContext(ctx, `CREATE TABLE IF NOT EXISTS jobs (
		id TEXT PRIMARY KEY,
		workload BLOB NOT NULL,
		phase TEXT NOT NULL,
		exit_code INTEGER
	)`); err != nil {
		return fmt.Errorf("initialize job schema: %w", err)
	}
	return nil
}

// rowReader is satisfied by both a database and an open write transaction.
type rowReader interface {
	QueryRowContext(context.Context, string, ...any) *sql.Row
}

// readJob decodes the latest committed or transaction-local record.
func readJob(ctx context.Context, reader rowReader, id string) (*Record, error) {
	record := &Record{ID: id, Workload: &jobpb.Workload{}}
	var workload []byte
	var exitCode sql.NullInt32
	err := reader.QueryRowContext(ctx, `SELECT workload, phase, exit_code FROM jobs WHERE id = ?`, id).
		Scan(&workload, &record.Phase, &exitCode)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("read job %s: %w", id, err)
	}
	if err := proto.Unmarshal(workload, record.Workload); err != nil {
		return nil, fmt.Errorf("decode workload for %s: %w", id, err)
	}
	if exitCode.Valid {
		record.ExitCode = &exitCode.Int32
	}
	return record, nil
}

// Get returns a committed SQLite record.
func (repository *Repository) Get(ctx context.Context, id string) (*Record, error) {
	record, err := readJob(ctx, repository.db, id)
	if err != nil {
		return nil, fmt.Errorf("read committed job: %w", err)
	}
	return record, nil
}

// Update loads a record and commits the changes requested by the callback.
// The callback indicates whether it modified the record.
func (repository *Repository) Update(ctx context.Context, id string, apply func(*Record, bool) (bool, error)) error {
	tx, err := repository.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin transaction for %s: %w", id, err)
	}
	defer tx.Rollback() // A committed transaction has nothing to roll back.
	record, err := readJob(ctx, tx, id)
	exists := err == nil
	if errors.Is(err, ErrNotFound) {
		record = &Record{ID: id}
	} else if err != nil {
		return fmt.Errorf("load transaction record: %w", err)
	}
	changed, err := apply(record, exists)
	if err != nil {
		return fmt.Errorf("apply job update: %w", err)
	}
	if changed {
		encoded, err := proto.Marshal(record.Workload)
		if err != nil {
			return fmt.Errorf("encode workload for %s: %w", id, err)
		}
		if _, err := tx.ExecContext(ctx, `INSERT INTO jobs(id, workload, phase, exit_code)
			VALUES (?, ?, ?, ?) ON CONFLICT(id) DO UPDATE SET
			phase = excluded.phase, exit_code = excluded.exit_code`,
			id, encoded, record.Phase, record.ExitCode); err != nil {
			return fmt.Errorf("persist job %s: %w", id, err)
		}
	}
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit job %s: %w", id, err)
	}
	return nil
}

// IDsInPhase lists job IDs in the supplied phase.
func (repository *Repository) IDsInPhase(ctx context.Context, phase string) (ids []string, err error) {
	rows, err := repository.db.QueryContext(ctx, `SELECT id FROM jobs WHERE phase = ?`, phase)
	if err != nil {
		return nil, fmt.Errorf("list worker jobs: %w", err)
	}
	defer func() { err = errors.Join(err, rows.Close()) }()
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return nil, fmt.Errorf("read active job ID: %w", err)
		}
		ids = append(ids, id)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate active worker jobs: %w", err)
	}
	return ids, nil
}
