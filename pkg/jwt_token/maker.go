package jwt_token

import (
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"job_search_platform/pkg/entities/common"
	"job_search_platform/pkg/helpers/server"

	"time"
)

const minSecretKeySize = 32

// Maker is an interface for managing tokens
type Maker interface {
	CreateToken(user common.UserResponse, tokenType string) (token string, payload *Payload, err error)
	VerifyToken(token string) (*Payload, error)
	GetErrorCode(err error) (errorCode int32)
}

// JWTMaker is a JSON Web Token maker
type JWTMaker struct {
	secretKey            string
	accessTokenDuration  time.Duration
	refreshTokenDuration time.Duration
}

func NewJWTMaker(secretKey string, accessTokenDuration time.Duration, refreshTokenDuration time.Duration) (Maker, error) {
	if len(secretKey) < minSecretKeySize {
		return nil, fmt.Errorf("invalid key size: must be at least %d characters", minSecretKeySize)
	}
	return &JWTMaker{
		secretKey: secretKey, accessTokenDuration: accessTokenDuration,
		refreshTokenDuration: refreshTokenDuration}, nil
}

func (maker *JWTMaker) CreateToken(user common.UserResponse, tokenType string) (token string, payload *Payload, err error) {
	switch tokenType {
	case "refresh":
		payload, err = NewPayload(user, tokenType, maker.refreshTokenDuration)
	case "access":
		payload, err = NewPayload(user, tokenType, maker.accessTokenDuration)
	}
	if err != nil {
		return "", payload, err
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	token, err = jwtToken.SignedString([]byte(maker.secretKey))
	return token, payload, err
}

func (maker *JWTMaker) VerifyToken(token string) (*Payload, error) {
	keyFunc := func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, ErrInvalidToken
		}
		return []byte(maker.secretKey), nil
	}

	jwtToken, err := jwt.ParseWithClaims(token, &Payload{}, keyFunc)
	if err != nil {
		var verr *jwt.ValidationError
		ok := errors.As(err, &verr)
		if ok && errors.Is(verr.Inner, ErrExpiredToken) {
			return nil, ErrExpiredToken
		}
		return nil, ErrInvalidToken
	}

	payload, ok := jwtToken.Claims.(*Payload)
	if !ok {
		return nil, ErrInvalidToken
	}

	return payload, nil
}

func (maker *JWTMaker) GetErrorCode(err error) (errorCode int32) {
	errorCode = -1
	if errors.Is(err, ErrExpiredToken) {
		return server.JWT_EXPIRES_ERR_CODE
	} else if errors.Is(err, ErrInvalidToken) {
		return server.TOKEN_VALIDATION_ERR_CODE
	}
	return errorCode
}
