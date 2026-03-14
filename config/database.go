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
	if err := ensureDatabaseExists(cfg); err != nil {
		return nil, err
	}

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

func ensureDatabaseExists(cfg *Config) error {
	postgresDSN := fmt.Sprintf("postgresql://%s:%s@%s:%s/postgres?sslmode=disable",
		cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort)

	ctx := context.Background()
	pool, err := pgxpool.New(ctx, postgresDSN)
	if err != nil {
		return fmt.Errorf("failed to connect to postgres database: %w", err)
	}
	defer pool.Close()

	var exists bool
	query := "SELECT EXISTS(SELECT 1 FROM pg_database WHERE datname = $1)"
	err = pool.QueryRow(ctx, query, cfg.DBName).Scan(&exists)
	if err != nil {
		return fmt.Errorf("failed to check if database exists: %w", err)
	}

	if !exists {
		dbLog.Infof("Database %s does not exist, creating it...", cfg.DBName)
		_, err = pool.Exec(ctx, fmt.Sprintf("CREATE DATABASE \"%s\"", cfg.DBName))
		if err != nil {
			return fmt.Errorf("failed to create database: %w", err)
		}
		dbLog.Infof("Database %s created successfully", cfg.DBName)
	}

	return nil
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
