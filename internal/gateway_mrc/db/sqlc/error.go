package db

import (
	"errors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"strconv"
)

const (
	ForeignKeyViolation = "23503"
	UniqueViolation     = "23505"
)

var ErrRecordNotFound = pgx.ErrNoRows

var ErrUniqueViolation = &pgconn.PgError{
	Code: UniqueViolation,
}

func ErrorCode(err error) int32 {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		i, _ := strconv.ParseInt(pgErr.Code, 10, 32)
		return int32(i)
	} else if errors.Is(err, pgx.ErrNoRows) {
		return -1
	}
	return 1
}
