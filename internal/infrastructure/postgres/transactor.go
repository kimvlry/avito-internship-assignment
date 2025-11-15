package postgres

import (
    "context"
    "fmt"
    "github.com/kimvlry/avito-internship-assignment/internal/domain/service"

    "github.com/jackc/pgx/v5"
    "github.com/jackc/pgx/v5/pgxpool"
)

type transactor struct {
    pool *pgxpool.Pool
}

func NewTransactor(pool *pgxpool.Pool) service.Transactor {
    return &transactor{pool: pool}
}

func (t *transactor) WithinTransaction(ctx context.Context, fn func(ctx context.Context) error) error {
    tx, err := t.pool.BeginTx(ctx, pgx.TxOptions{
        IsoLevel:   pgx.ReadCommitted,
        AccessMode: pgx.ReadWrite,
    })
    if err != nil {
        return fmt.Errorf("begin transaction: %w", err)
    }

    txCtx := injectTx(ctx, tx)

    defer func() {
        if p := recover(); p != nil {
            _ = tx.Rollback(ctx)
            panic(p)
        }
    }()

    if err := fn(txCtx); err != nil {
        if rbErr := tx.Rollback(ctx); rbErr != nil {
            return fmt.Errorf("rollback transaction (original error: %w): %v", err, rbErr)
        }
        return err
    }

    if err := tx.Commit(ctx); err != nil {
        return fmt.Errorf("commit transaction: %w", err)
    }
    return nil
}
