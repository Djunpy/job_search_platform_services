package db

import (
	"context"
	"github.com/jackc/pgx/v5/pgtype"
	"job_search_platform/pkg/helpers/crypto"
)

type UserPhone struct {
	Number      int64  `json:"number"`
	CountryCode string `json:"country_code"`
}

type BaseUserInfo struct {
	Email     string `json:"email" validate:"required,email"`
	Password1 string `json:"password1" validate:"required,password"`
	Password2 string `json:"password2" validate:"required,eqfield=Password1"`
	FirstName string `json:"first_name" validate:"required"`
	LastName  string `json:"last_name" validate:"required"`
	Source    string `json:"source,omitempty"`
	UserType  string `json:"user_type" validate:"required"` // 'company', 'job_seeker'
	UserSexy  string `json:"user_sexy"`
}

//Google Auth → google_auth
//Apple Auth →  apple_auth
//Обычная регистрация → standard_auth
//Админ-панель → admin_auth

type CreateOrdinaryUserTxParams struct {
	BaseUserInfo
	UserPhone `json:"phone" validate:"required"`
}

type CreateStaffUserTxParams struct {
	BaseUserInfo
	InviteCode string `json:"invite,omitempty"`
}

func (store *SQLStore) TxCreateUser(ctx context.Context, args *CreateOrdinaryUserTxParams) error {
	err := store.execTx(ctx, func(q *Queries) error {
		var err error
		userType := UserTypes(args.UserType)
		userSexy := Sexy(args.UserSexy)
		hashPass := crypto.HashPassword(args.Password1)
		userArgs := &CreateUserParams{
			Email:      args.Email,
			FirstName:  pgtype.Text{String: args.FirstName, Valid: args.FirstName != ""},
			LastName:   pgtype.Text{String: args.LastName, Valid: args.LastName != ""},
			Password:   hashPass,
			AuthSource: args.Source,
			UserType:   NullUserTypes{UserTypes: userType, Valid: true},
			Sexy:       NullSexy{Sexy: userSexy, Valid: userSexy != ""},
		}
		group, err := store.GetGroupByName(ctx, "ordinary_users")
		if err != nil {
			return err
		}

		user, err := q.CreateUser(ctx, *userArgs)
		if err != nil {
			return err
		}

		phoneArgs := &CreateUserPhoneParams{
			UserID:      user.ID,
			Number:      args.Number,
			CountryCode: args.CountryCode,
		}
		_, err = q.CreateUserPhone(ctx, *phoneArgs)
		if err != nil {
			return err
		}

		groupArgs := &CreateUserGroupParams{
			GroupID: group.ID,
			UserID:  user.ID,
		}
		_, err = q.CreateUserGroup(ctx, *groupArgs)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}
