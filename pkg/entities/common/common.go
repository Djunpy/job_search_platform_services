package common

import "github.com/google/uuid"

type CommonResponse struct {
	Code  int    `json:"code"`
	Error string `json:"error"`
}

type UserResponse struct {
	UserId   uuid.UUID `json:"user_id"`
	Email    string    `json:"email"`
	Groups   []string  `json:"roles"`
	UserType string    `json:"user_type"`
}

type SignInBodyResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type SignInResponse struct {
	CommonResponse
	Body SignInBodyResponse `json:"body"`
}

type RefreshTokenBodyResponse struct {
	AccessToken string `json:"access_token"`
}

type RefreshTokenResponse struct {
	CommonResponse
	RefreshTokenBodyResponse `json:"body"`
}

type PayloadSendVerifyEmail struct {
	Email     string `json:"email"`
	JWTToken  string `json:"jwt_token"`
	LangCode  string `json:"lang_code"`
	LastName  string `json:"last_name"`
	FirstName string `json:"first_name"`
}
