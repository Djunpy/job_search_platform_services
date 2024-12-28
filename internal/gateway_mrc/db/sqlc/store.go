package db

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Store interface {
	Querier
	GetOrCreateClientSession(
		ctx context.Context, sessionId string, req RequestArgs) (session Session, created bool, errCode int32, err error)
	UpdateSessionLastActive(ctx context.Context, sessionId string) (session Session, err error)
}

type SQLStore struct {
	connPool *pgxpool.Pool
	*Queries
}

func NewStore(connPool *pgxpool.Pool) *SQLStore {
	return &SQLStore{
		connPool: connPool,
		Queries:  New(connPool),
	}
}
