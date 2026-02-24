package config

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/op/go-logging"
)

var dbLog = logging.MustGetLogger("database")

type Database struct {
	Pool *pgxpool.Pool
}

func NewDatabase(cfg *Config) (*Database, error) {
	dsn := cfg.GetDSN()
	config, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to parse database config: %w", err)
	}

	pool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		return nil, fmt.Errorf("failed to create connection pool: %w", err)
	}

	if err := pool.Ping(context.Background()); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	dbLog.Info("Database connected successfully")
	return &Database{Pool: pool}, nil
}

func (db *Database) InitSchema(schemaPath string) error {
	schema, err := os.ReadFile(schemaPath)
	if err != nil {
		return fmt.Errorf("failed to read schema file: %w", err)
	}

	_, err = db.Pool.Exec(context.Background(), string(schema))
	if err != nil {
		return fmt.Errorf("failed to execute schema: %w", err)
	}

	dbLog.Info("Database schema initialized successfully")
	return nil
}

func (db *Database) Close() {
	db.Pool.Close()
}
