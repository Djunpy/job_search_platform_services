package jwt_token

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"job_search_platform/pkg/entities/common"
	"time"
)

var (
	ErrInvalidToken = errors.New("token is invalid")
	ErrExpiredToken = errors.New("token has expired")
)

type Payload struct {
	TokenType string    `json:"token_type"`
	ID        uuid.UUID `json:"id"`
	UserId    uuid.UUID `json:"user_id"`
	//Username    string    `json:"username"`
	Email    string `json:"email"`
	IsActive bool   `json:"is_active"`
	//IsSuperuser bool      `json:"is_superuser"`
	//IsStaff     bool      `json:"is_staff"`
	Groups    []string  `json:"roles"`
	IssuedAt  time.Time `json:"issued_at"`
	ExpiredAt time.Time `json:"expired_at"`
}

func NewPayload(user common.UserResponse, tokenType string, duration time.Duration) (*Payload, error) {
	tokenID, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}

	payload := &Payload{
		TokenType: tokenType,
		ID:        tokenID,
		UserId:    user.UserId,
		//Username:    user.Username,
		Email: user.Email,
		//IsActive: user.IsActive,
		//IsSuperuser: user.IsSuperuser,
		//IsStaff:     user.IsStaff,
		Groups:    user.Groups,
		IssuedAt:  time.Now(),
		ExpiredAt: time.Now().Add(duration),
	}
	return payload, nil
}

func (payload *Payload) Valid() error {
	if time.Now().After(payload.ExpiredAt) {
		return ErrExpiredToken
	}
	return nil
}

func GetJWTPayload(ctx *gin.Context) (*Payload, bool) {
	var jwtPayload *Payload
	ctxPayload, exists := ctx.Get("jwtTokenPayload")
	if !exists {
		return jwtPayload, false
	}

	jwtPayload, ok := ctxPayload.(*Payload)
	if !ok {
		return jwtPayload, false
	}

	return jwtPayload, true
}
