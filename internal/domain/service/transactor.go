package service

import "context"

type Transactor interface {
    WithinTransaction(ctx context.Context, tFunc func(ctx context.Context) error) error
}
