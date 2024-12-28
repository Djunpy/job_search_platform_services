package usecases

import (
	"context"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	db "job_search_platform/internal/users_mrc/db/sqlc"
	"job_search_platform/internal/users_mrc/entities"
	"job_search_platform/pkg/database"
	"job_search_platform/pkg/helpers/data_processing"
	"job_search_platform/pkg/helpers/server"
	"job_search_platform/pkg/jwt_token"
)

type UsersUsecase struct {
	store db.Store
}

func NewUsersUsecase(store db.Store) UsersUsecase {
	return UsersUsecase{store}
}

func (usecase *UsersUsecase) GetUserDetail(
	ctx context.Context, jwtPayload jwt_token.Payload) (accountDetail entities.UserDetail, statusCode int32, err error) {
	user, err := usecase.store.GetUserAndGroupsByEmail(ctx, jwtPayload.Email)
	if err != nil {
		return accountDetail, database.ErrorCode(err), err
	}
	var groups []string
	groups, err = data_processing.ExtractGroups(user.Groups)
	if err != nil {
		return accountDetail, database.ErrorCode(err), err

	}
	roles := data_processing.RemoveSAtEnd(groups)
	phone, err := usecase.store.GetUserPhoneByUserId(ctx, user.ID)
	if err != nil {
		return accountDetail, database.ErrorCode(err), err
	}
	accountDetail = entities.NewUserDetailResponse(user, phone, roles)
	return accountDetail, server.SUCCESS_CODE, nil
}

func (usecase *UsersUsecase) UpdateUser(ctx context.Context, payload entities.UserUpdate, userId uuid.UUID) (statusCode int32, err error) {
	if payload.FirstName != "" || payload.LastName != "" {
		updateUserParams := db.UpdateUserByIdParams{
			FirstName: pgtype.Text{String: payload.FirstName, Valid: payload.FirstName != ""},
			LastName:  pgtype.Text{String: payload.LastName, Valid: payload.LastName != ""},
		}
		_, err = usecase.store.UpdateUserById(ctx, updateUserParams)
		if err != nil {
			return database.ErrorCode(err), err
		}
	}
	if payload.Phone != 0 || payload.CountryCode != "" {
		updatePhoneParams := db.UpdateUserPhoneParams{
			Number:      pgtype.Int8{Int64: payload.Phone, Valid: payload.Phone != 0},
			CountryCode: pgtype.Text{String: payload.CountryCode, Valid: payload.CountryCode != ""},
		}
		_, err = usecase.store.UpdateUserPhone(ctx, updatePhoneParams)
	}
	return server.SUCCESS_CODE, nil
}
