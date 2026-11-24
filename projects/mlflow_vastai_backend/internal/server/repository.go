package server

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"git.alwaldend.com/alwaldend/src/projects/mlflow_vastai_backend/api/jobpb"
	"google.golang.org/protobuf/proto"
)

// assignment stores the request, chosen rental, and latest confirmed job.
type assignment struct {
	// submission is immutable after reservation.
	submission *jobpb.SubmitJobRequest
	// instance is the assigned Vast.ai contract ID.
	instance int64
	// job is the last confirmed worker observation.
	job *jobpb.Job
	// cancelling records an outstanding cancellation during transport failures.
	cancelling bool
	// settled indicates that reconciliation has no remaining work.
	settled bool
}

// repository persists the coordinator's assignments using a borrowed connection.
type repository struct {
	// db stores assignments in the coordinator database.
	db *sql.DB
}

// openRepository ensures the coordinator schema before accepting jobs.
func openRepository(ctx context.Context, db *sql.DB) (*repository, error) {
	_, err := db.ExecContext(ctx, `CREATE TABLE IF NOT EXISTS jobs (
		id TEXT PRIMARY KEY,
		submission BLOB NOT NULL,
		instance INTEGER NOT NULL,
		observation BLOB NOT NULL,
		cancelling BOOLEAN NOT NULL DEFAULT FALSE,
		settled BOOLEAN NOT NULL DEFAULT FALSE
	)`)
	if err != nil {
		return nil, fmt.Errorf("initialize assignment schema: %w", err)
	}
	return &repository{db: db}, nil
}

// find reads a job and the assignment needed to query its worker.
func (repository *repository) find(ctx context.Context, id string) (*assignment, error) {
	row := &assignment{submission: &jobpb.SubmitJobRequest{}, job: &jobpb.Job{}}
	var submission, observation []byte
	err := repository.db.QueryRowContext(ctx, `SELECT submission, instance, observation, cancelling, settled
		FROM jobs WHERE id = ?`, id).Scan(&submission, &row.instance, &observation, &row.cancelling, &row.settled)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("read assignment: %w", err)
	}
	if err := proto.Unmarshal(submission, row.submission); err != nil {
		return nil, fmt.Errorf("decode stored submission: %w", err)
	}
	if err := proto.Unmarshal(observation, row.job); err != nil {
		return nil, fmt.Errorf("decode stored job: %w", err)
	}
	return row, nil
}

// reserve persists the assigned rental before any SSH command is issued.
func (repository *repository) reserve(ctx context.Context, row *assignment) error {
	submission, err := proto.Marshal(row.submission)
	if err != nil {
		return fmt.Errorf("encode submission reservation: %w", err)
	}
	observation, err := proto.Marshal(row.job)
	if err != nil {
		return fmt.Errorf("encode initial result: %w", err)
	}
	_, err = repository.db.ExecContext(ctx,
		`INSERT INTO jobs (id, submission, instance, observation) VALUES (?, ?, ?, ?)`,
		row.job.GetJobId(), submission, row.instance, observation)
	if err != nil {
		return fmt.Errorf("reserve assignment: %w", err)
	}
	return nil
}

// save commits the latest observation and remaining reconciliation work.
func (repository *repository) save(ctx context.Context, row *assignment) error {
	observation, err := proto.Marshal(row.job)
	if err != nil {
		return fmt.Errorf("encode job observation: %w", err)
	}
	_, err = repository.db.ExecContext(ctx, `UPDATE jobs SET observation = ?, cancelling = ?, settled = ? WHERE id = ?`,
		observation, row.cancelling, row.settled, row.job.GetJobId())
	if err != nil {
		return fmt.Errorf("save job observation: %w", err)
	}
	return nil
}

// pending lists assignments that still require a worker or tracking observation.
func (repository *repository) pending(ctx context.Context) (_ []string, err error) {
	rows, err := repository.db.QueryContext(ctx, `SELECT id FROM jobs WHERE settled = FALSE ORDER BY id`)
	if err != nil {
		return nil, fmt.Errorf("list pending assignments: %w", err)
	}
	defer func() { err = errors.Join(err, rows.Close()) }()
	var ids []string
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return nil, fmt.Errorf("read pending job ID: %w", err)
		}
		ids = append(ids, id)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate pending jobs: %w", err)
	}
	return ids, nil
}
