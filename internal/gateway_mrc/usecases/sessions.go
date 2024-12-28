package usecases

import (
	"context"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	db "job_search_platform/internal/gateway_mrc/db/sqlc"
	"job_search_platform/pkg/helpers/server"
	"time"
)

type SessionsUsecase struct {
	store db.Store
}

func NewSessionsUsecase(store db.Store) SessionsUsecase {
	return SessionsUsecase{store}
}

func (uc *SessionsUsecase) UpdateSession(ctx context.Context, sessionIdStr string, refreshToken, accessToken string) (db.Session, int32, error) {
	var session db.Session
	var err error
	sessionId, err := uuid.Parse(sessionIdStr)
	if err != nil {
		return session, server.SESSION_PARSING_ERR_CODE, err
	}
	sessionArgs := &db.UpdateSessionDataParams{
		AccessToken:  pgtype.Text{String: accessToken, Valid: accessToken != ""},
		RefreshToken: pgtype.Text{String: refreshToken, Valid: refreshToken != ""},
		ID:           pgtype.UUID{Bytes: sessionId, Valid: true},
	}
	session, err = uc.store.UpdateSessionData(ctx, *sessionArgs)
	if err != nil {
		return session, db.ErrorCode(err), err
	}
	return session, server.SUCCESS_CODE, nil
}

func (uc *SessionsUsecase) UpdateSessionLastActive(ctx context.Context, sessionId uuid.UUID) (db.Session, int32, error) {
	sessionArgs := &db.UpdateSessionDataParams{
		LastActive: pgtype.Timestamptz{Time: time.Now(), Valid: true},
		ID:         pgtype.UUID{Bytes: sessionId, Valid: true},
	}
	session, err := uc.store.UpdateSessionData(ctx, *sessionArgs)
	if err != nil {
		return session, db.ErrorCode(err), err
	}
	return session, server.SUCCESS_CODE, nil
}

func (uc *SessionsUsecase) GetSession(ctx context.Context, id string) (db.Session, int32, error) {
	var session db.Session
	sessionId, err := uuid.Parse(id)
	if err != nil {
		return session, db.ErrorCode(err), err
	}
	session, err = uc.store.GetSession(ctx, pgtype.UUID{Bytes: sessionId, Valid: true})
	if err != nil {
		return session, db.ErrorCode(err), err
	}
	return session, db.ErrorCode(err), nil
}
