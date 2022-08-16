package postgres

import (
	"context"
	"fmt"
	"log"

	"gitlab.ozon.dev/igor.benko.1991/homework/internal/config"
	"gitlab.ozon.dev/igor.benko.1991/homework/migrations"

	"github.com/jackc/pgx/v4/pgxpool"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/pressly/goose/v3"
)

func New(ctx context.Context, cfg *config.PoolerConfig) (*pgxpool.Pool, error) {
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Name)

	pool, err := pgxpool.Connect(ctx, dsn)
	if err != nil {
		log.Fatal("can't connect to database", err)
	}

	if err := pool.Ping(ctx); err != nil {
		log.Fatal("ping database error", err)
	}

	config := pool.Config()
	config.ConnConfig.PreferSimpleProtocol = true

	return pool, nil
}

func Migrate(ctx context.Context, cfg *config.DatabaseConfig) error {
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Name)

	db, err := sqlx.Connect("pgx", dsn)
	if err != nil {
		return err
	}

	defer db.Close()

	goose.SetBaseFS(migrations.MigrationFiles)

	if err := goose.Up(db.DB, "."); err != nil {
		return err
	}

	return nil
}
