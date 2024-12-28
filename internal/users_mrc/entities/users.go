package entities

import (
	"github.com/google/uuid"
	db "job_search_platform/internal/users_mrc/db/sqlc"
)

type UserUpdate struct {
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	Phone       int64  `json:"phone"`
	CountryCode string `json:"country_code"`
}

type UserPhone struct {
	Number      int64  `json:"number"`
	CountryCode string `json:"country_code"`
}

type UserDetail struct {
	Id            uuid.UUID `json:"user_id"`
	VerifiedEmail bool      `json:"verified_email"`
	Email         string    `json:"email"`
	Roles         []string  `json:"roles"`
	UserType      string    `json:"user_type"`
	FirstName     string    `json:"first_name"`
	LastName      string    `json:"last_name"`
	AccountPhone  UserPhone `json:"phone"`
}

func NewUserDetailResponse(user db.GetUserAndGroupsByEmailRow, phone db.Phone, roles []string) UserDetail {
	userType := user.UserType
	account := UserDetail{
		Id:            user.ID.Bytes,
		Email:         user.Email,
		VerifiedEmail: user.VerifiedEmail.Bool,
		Roles:         roles,
		UserType:      string(userType.UserTypes),
		FirstName:     user.FirstName.String,
		LastName:      user.LastName.String,
	}
	account.AccountPhone.Number = phone.Number
	account.AccountPhone.CountryCode = phone.CountryCode
	return account
}
