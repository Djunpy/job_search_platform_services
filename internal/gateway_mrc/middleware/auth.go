package middleware

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"io"
	db "job_search_platform/internal/gateway_mrc/db/sqlc"
	"job_search_platform/pkg/config"
	"job_search_platform/pkg/entities/common"
	"job_search_platform/pkg/helpers/server"
	"job_search_platform/pkg/jwt_token"
	"net/http"
)

func AuthMiddleware(tokenMaker jwt_token.Maker, store db.Store, config config.Config) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var statusCode int32
		var jwtPayload *jwt_token.Payload
		sessionIdStr, err := ctx.Cookie("session_id")
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, server.Response(err, server.GET_COOKIE_ERR_CODE, nil))
			return
		}

		sessionId, err := uuid.Parse(sessionIdStr)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, server.Response(err, server.GET_COOKIE_ERR_CODE, nil))
			return
		}

		session, err := store.GetSession(ctx, pgtype.UUID{Bytes: sessionId, Valid: true})
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, server.Response(err, server.SESSION_NOT_FOUND_ERR_CODE, nil))
			return
		}

		if session.IsBlocked.Bool {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, server.Response(err, server.SESSION_BLOCKED_ERR_CODE, nil))
			return
		}

		_, err = tokenMaker.VerifyToken(session.RefreshToken.String)
		if err != nil {
			statusCode = tokenMaker.GetErrorCode(err)
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, server.Response(err, statusCode, nil))
			return
		}
		jwtPayload, err = tokenMaker.VerifyToken(session.AccessToken.String)
		if err != nil {
			statusCode = tokenMaker.GetErrorCode(err)
			if statusCode == server.JWT_EXPIRES_ERR_CODE {
				tokenRefreshEndpoint := fmt.Sprintf("%s/%s", config.UsersMrcUrl, "api/v1/auth/public/refresh-token")
				newAccessToken, refreshErr := refreshAccessToken(session.RefreshToken.String, tokenRefreshEndpoint)
				if refreshErr != nil {
					ctx.AbortWithStatusJSON(http.StatusUnauthorized, server.Response(refreshErr, server.SENDING_TOKEN_REFRESH_ERR_CODE, nil))
					return
				}
				sessionArgs := &db.UpdateSessionDataParams{
					AccessToken: pgtype.Text{String: newAccessToken, Valid: newAccessToken != ""},
					ID:          pgtype.UUID{Bytes: sessionId, Valid: true},
				}
				_, err := store.UpdateSessionData(ctx, *sessionArgs)
				if refreshErr != nil {
					ctx.AbortWithStatusJSON(http.StatusUnauthorized, server.Response(err, server.SENDING_TOKEN_REFRESH_ERR_CODE, nil))
					return
				}
				jwtPayload, _ = tokenMaker.VerifyToken(session.AccessToken.String)
			} else {
				ctx.AbortWithStatusJSON(http.StatusUnauthorized, server.Response(err, statusCode, nil))
				return
			}
		}
		ctx.Set("jwtTokenPayload", jwtPayload)
		ctx.Next()
	}
}

func refreshAccessToken(refreshToken string, endpoint string) (string, error) {
	// Создаем запрос на обновление токена
	var respData common.RefreshTokenResponse
	reqBody := map[string]string{"refresh_token": refreshToken}
	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return "", err
	}
	resp, err := http.Post(endpoint, "application/json", bytes.NewBuffer(jsonBody))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to refresh token, status code: %d", resp.StatusCode)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	err = json.Unmarshal(body, &respData)
	if err != nil {
		return "", err
	}
	return respData.AccessToken, nil
}
