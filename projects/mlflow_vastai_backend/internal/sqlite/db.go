package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/url"

	"git.alwaldend.com/alwaldend/src/projects/mlflow_vastai_backend/internal/machine"
	_ "modernc.org/sqlite" // Registers the database/sql driver.
)

const (
	// stateClosed has no connection available to borrow.
	stateClosed machine.State = "sqlite-closed"
	// stateOpen holds a live connection.
	stateOpen machine.State = "sqlite-open"
	// open connects to the SQLite file.
	open machine.Event = "open"
	// close releases the SQLite connection.
	close machine.Event = "close"
)

// DB owns a SQLite connection. Its owner closes it after all borrowers finish.
type DB struct {
	// definition controls the connection lifetime.
	definition *machine.Machine[string]
	// state records the connection phase.
	state machine.State
	// connection is borrowed by callers for queries and transactions.
	connection *sql.DB
}

// newDB configures a closed connection owner without touching the filesystem.
func newDB() *DB {
	db := &DB{definition: machine.New[string](stateClosed), state: stateClosed}
	db.definition.Configure(stateClosed).Permit(open, stateOpen).Ignore(close).
		OnEntryFrom(close, db.disconnect)
	db.definition.Configure(stateOpen).OnEntryFrom(open, db.connect).Permit(close, stateClosed)
	return db
}

// Open connects to the SQLite file at the supplied absolute path.
func Open(ctx context.Context, path string) (*DB, error) {
	db := newDB()
	ctx = db.definition.WithInvocation(ctx, &db.state, path)
	if err := db.definition.FireCtx(ctx, open); err != nil {
		return nil, fmt.Errorf("open SQLite %q: %w", path, err)
	}
	return db, nil
}

// Connection returns the handle borrowed for the open lifetime of DB.
func (db *DB) Connection() *sql.DB {
	return db.connection
}

// Close releases the connection after its borrowers have stopped.
func (db *DB) Close() error {
	ctx := db.definition.WithInvocation(context.Background(), &db.state, "")
	if err := db.definition.FireCtx(ctx, close); err != nil {
		return fmt.Errorf("close SQLite: %w", err)
	}
	return nil
}

// connect establishes the connection before making it available to callers.
func (db *DB) connect(ctx context.Context, _ ...any) error {
	path, err := db.definition.Input(ctx)
	if err != nil {
		return fmt.Errorf("read SQLite path: %w", err)
	}
	dsn := url.URL{Scheme: "file", Path: path}
	dsn.RawQuery = "_pragma=busy_timeout(5000)&_pragma=journal_mode(WAL)&_txlock=immediate"
	connection, err := sql.Open("sqlite", dsn.String())
	if err != nil {
		return fmt.Errorf("create connection: %w", err)
	}
	connection.SetMaxOpenConns(1)
	if err := connection.PingContext(ctx); err != nil {
		return errors.Join(fmt.Errorf("connect to SQLite: %w", err), connection.Close())
	}
	db.connection = connection
	return nil
}

// disconnect closes the connection as the transition returns to closed.
func (db *DB) disconnect(_ context.Context, _ ...any) error {
	connection := db.connection
	db.connection = nil
	if err := connection.Close(); err != nil {
		return fmt.Errorf("release SQLite connection: %w", err)
	}
	return nil
}

// Graph describes the database lifecycle without opening a file.
func Graph() machine.GraphSection {
	return machine.GraphSection{ID: "sqlite", Title: "SQLite connection", DOT: newDB().definition.ToGraph()}
}
