package postgres

import (
	"context"
	"embed"
	"fmt"
	"io/fs"
	"time"

	"github.com/Neeeooshka/gopher-club/internal/storage/postgres/sqlc"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/tern/v2/migrate"
)

const versionTable = "gophermart_schema_versions"

//go:embed migrations/*.sql
var migrationFiles embed.FS

type Postgres struct {
	DB   *pgxpool.Pool
	sqlc *sqlc.Queries
}

func (s *Postgres) Close() error {
	s.DB.Close()
	return nil
}

func NewPostgresStorage(conn string) (pgx *Postgres, err error) {

	pgx = &Postgres{}

	cfg, err := pgxpool.ParseConfig(conn)
	if err != nil {
		return nil, fmt.Errorf("error parsing connection string: %w", err)
	}

	pgx.DB, err = pgxpool.NewWithConfig(context.Background(), cfg)
	if err != nil {
		return nil, fmt.Errorf("error connecting to database: %w", err)
	}

	pgx.sqlc = sqlc.New(pgx.DB)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	err = pgx.runMigrations(ctx)
	if err != nil {
		return nil, fmt.Errorf("error applying migrations: %w", err)
	}

	return pgx, nil
}

func (s *Postgres) Ping() error {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return s.DB.Ping(ctx)
}

func (s *Postgres) runMigrations(ctx context.Context) error {

	conn, err := s.DB.Acquire(ctx)
	if err != nil {
		return fmt.Errorf("error DB connection: %w", err)
	}
	defer conn.Release()

	migrator, err := migrate.NewMigrator(ctx, conn.Conn(), versionTable)
	if err != nil {
		return fmt.Errorf("cannot construct migrator: %w", err)
	}

	migrationRoot, err := fs.Sub(migrationFiles, "migrations")
	if err != nil {
		return fmt.Errorf("error loading migration root: %w", err)
	}

	if err := migrator.LoadMigrations(migrationRoot); err != nil {
		return fmt.Errorf("error loading migrations: %w", err)
	}

	return migrator.Migrate(ctx)
}
