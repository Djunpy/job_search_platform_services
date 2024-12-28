package entities

type SignInReq struct {
	Login    string `json:"login" validate:"required,min=6"`
	Password string `json:"password" validate:"required,min=6"`
}

type RefreshTokenReq struct {
	RefreshToken string `json:"refresh_token"`
}

type ChangePasswordReq struct {
	OldPassword string `json:"old_Password" validate:"required,password"`
	NewPassword string `json:"new_password" validate:"required,password"`
}
