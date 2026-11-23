package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/url"
	"path/filepath"
	"strings"

	"git.alwaldend.com/alwaldend/src/projects/mlflow_vastai_backend/api/jobpb"
	"google.golang.org/protobuf/proto"
	_ "modernc.org/sqlite" // Registers the SQLite database/sql driver.
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

// Store manages the worker's SQLite database.
type Store struct {
	// db serializes local access to the job table.
	db *sql.DB
}

// Open configures SQLite for concurrent processes and initializes its schema.
func Open(ctx context.Context, directory string) (*Store, error) {
	dsn := url.URL{Scheme: "file", Path: filepath.Join(directory, "jobs.sqlite3")}
	dsn.RawQuery = "_pragma=busy_timeout(5000)&_pragma=journal_mode(WAL)&_txlock=immediate"
	db, err := sql.Open("sqlite", dsn.String())
	if err != nil {
		return nil, fmt.Errorf("open worker SQLite: %w", err)
	}
	db.SetMaxOpenConns(1)
	if _, err := db.ExecContext(ctx, `CREATE TABLE IF NOT EXISTS jobs (
		id TEXT PRIMARY KEY,
		workload BLOB NOT NULL,
		phase TEXT NOT NULL,
		exit_code INTEGER
	)`); err != nil {
		return nil, errors.Join(fmt.Errorf("initialize worker SQLite: %w", err), db.Close())
	}
	return &Store{db: db}, nil
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
func (s *Store) Get(ctx context.Context, id string) (*Record, error) {
	record, err := readJob(ctx, s.db, id)
	if err != nil {
		return nil, fmt.Errorf("read committed job: %w", err)
	}
	return record, nil
}

// Update applies a callback inside a transaction and persists phase changes.
// The callback receives whether the record existed before the transaction.
func (s *Store) Update(ctx context.Context, id string, apply func(*Record, bool) error) error {
	tx, err := s.db.BeginTx(ctx, nil)
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
	before := record.Phase
	if err := apply(record, exists); err != nil {
		return fmt.Errorf("apply job update: %w", err)
	}
	if record.Phase != before {
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

// Phase reads the phase without decoding the workload.
func (s *Store) Phase(ctx context.Context, id string) (string, error) {
	var phase string
	err := s.db.QueryRowContext(ctx, `SELECT phase FROM jobs WHERE id = ?`, id).Scan(&phase)
	if errors.Is(err, sql.ErrNoRows) {
		return "", ErrNotFound
	}
	if err != nil {
		return "", fmt.Errorf("read phase for %s: %w", id, err)
	}
	return phase, nil
}

// IDsInPhases lists job IDs with any of the supplied phases.
func (s *Store) IDsInPhases(ctx context.Context, phases ...string) (ids []string, err error) {
	if len(phases) == 0 {
		return nil, nil
	}
	args := make([]any, len(phases))
	for i, phase := range phases {
		args[i] = phase
	}
	query := `SELECT id FROM jobs WHERE phase IN (` + strings.Repeat("?, ", len(phases)-1) + "?)"
	rows, err := s.db.QueryContext(ctx, query, args...)
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

// Close releases the listener's SQLite connection.
func (s *Store) Close() error {
	if err := s.db.Close(); err != nil {
		return fmt.Errorf("close worker SQLite: %w", err)
	}
	return nil
}
