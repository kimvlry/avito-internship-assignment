package postgres

import (
    "context"
    "fmt"
    "github.com/jackc/pgx/v5/pgconn"

    "github.com/Masterminds/squirrel"
    "github.com/jackc/pgx/v5"
    "github.com/jackc/pgx/v5/pgxpool"
)

type DB struct {
    *pgxpool.Pool
    qb squirrel.StatementBuilderType
}

func NewDB(ctx context.Context, connString string) (*DB, error) {
    config, err := pgxpool.ParseConfig(connString)
    if err != nil {
        return nil, fmt.Errorf("parse config: %w", err)
    }

    config.MaxConns = 25
    config.MinConns = 5

    pool, err := pgxpool.NewWithConfig(ctx, config)
    if err != nil {
        return nil, fmt.Errorf("create pool: %w", err)
    }

    if err := pool.Ping(ctx); err != nil {
        return nil, fmt.Errorf("ping database: %w", err)
    }

    return &DB{
        Pool: pool,
        qb:   squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
    }, nil
}

func (db *DB) Close() {
    db.Pool.Close()
}

func (db *DB) QueryBuilder() squirrel.StatementBuilderType {
    return db.qb
}

type txKeyType struct{}

var txKey = txKeyType{}

func injectTx(ctx context.Context, tx pgx.Tx) context.Context {
    return context.WithValue(ctx, txKey, tx)
}

func extractTx(ctx context.Context) (pgx.Tx, bool) {
    tx, ok := ctx.Value(txKey).(pgx.Tx)
    return tx, ok
}

func (db *DB) GetQuerier(ctx context.Context) Querier {
    if tx, ok := extractTx(ctx); ok {
        return tx
    }
    return db.Pool
}

type Querier interface {
    Exec(ctx context.Context, sql string, arguments ...interface{}) (pgconn.CommandTag, error)
    Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error)
    QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row
}
