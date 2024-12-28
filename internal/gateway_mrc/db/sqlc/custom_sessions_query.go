package db

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"

	"job_search_platform/pkg/config"
	"job_search_platform/pkg/helpers/server"
	"log"
	"time"
)

type RequestArgs struct {
	IpAddress string
	UserAgent string
}

func (store *SQLStore) CreateClientSession(ctx context.Context, clientId, userAgent string, sessionDuration int) (Session, error) {

	now := time.Now()

	thirtyDaysAhead := now.AddDate(0, 0, sessionDuration)
	sessionArgs := &CreateSessionParams{
		UserAgent: userAgent,
		ClientIp:  clientId,
		ExpiresAt: thirtyDaysAhead,
	}
	session, err := store.CreateSession(ctx, *sessionArgs)
	if err != nil {
		return session, err
	}
	return session, nil
}

func (store *SQLStore) GetOrCreateClientSession(
	ctx context.Context, sessionId string, req RequestArgs) (session Session, created bool, errCode int32, err error) {
	conf, err := config.LoadConfig(".", "gateway_mrc")
	if err != nil {
		log.Fatalf("could not load config: %v", err)
	}

	if sessionId != "" {
		parsedUUID, err := uuid.Parse(sessionId)
		if err != nil {
			return session, false, ErrorCode(err), err
		}
		session, err = store.GetSession(ctx, pgtype.UUID{Bytes: parsedUUID, Valid: true})
		if errors.Is(err, pgx.ErrNoRows) {
			session, err = store.CreateClientSession(ctx, req.IpAddress, req.UserAgent, conf.SessionDuration)
			return session, true, server.SUCCESS_CODE, nil
		}

		return session, false, server.SUCCESS_CODE, nil
	}
	session, err = store.CreateClientSession(ctx, req.IpAddress, req.UserAgent, conf.SessionDuration)
	if err != nil {
		return session, false, ErrorCode(err), err
	}
	return session, true, server.SUCCESS_CODE, nil
}

func (store *SQLStore) UpdateSessionLastActive(ctx context.Context, sessionId string) (session Session, err error) {
	var parsedUUID uuid.UUID
	if sessionId != "" {
		parsedUUID, err = uuid.Parse(sessionId)
		if err != nil {
			return session, err
		}
	}
	updateArgs := &UpdateSessionDataParams{
		ID:         pgtype.UUID{Bytes: parsedUUID, Valid: true},
		LastActive: pgtype.Timestamptz{Time: time.Now(), Valid: true},
	}
	session, err = store.UpdateSessionData(ctx, *updateArgs)
	if err != nil {
		return session, err
	}
	return session, nil
}
