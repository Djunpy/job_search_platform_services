package handlers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/hibiken/asynq"
	"github.com/rs/zerolog/log"
	db "job_search_platform/internal/users_mrc/db/sqlc"
	"job_search_platform/internal/users_mrc/entities"
	"job_search_platform/internal/users_mrc/usecases"
	"job_search_platform/pkg/entities/common"
	"job_search_platform/pkg/helpers/server"
	"job_search_platform/pkg/jwt_token"
	"job_search_platform/pkg/scheduler"
	"net/http"
	"time"
)

type AuthHandler struct {
	usecase         usecases.AuthUsecase
	taskDistributor scheduler.TaskDistributor
}

func NewAuthHandler(usecase usecases.AuthUsecase, taskDistributor scheduler.TaskDistributor) AuthHandler {
	return AuthHandler{usecase: usecase, taskDistributor: taskDistributor}
}

func (handler *AuthHandler) SignUpUser(ctx *gin.Context) {
	var payload *db.CreateOrdinaryUserTxParams
	var err error
	if err = ctx.ShouldBindJSON(&payload); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, server.Response(err, server.INVALID_DATA_ERR_CODE, nil))
		return
	}
	token, errCode, err := handler.usecase.CreateUser(
		ctx,
		payload,
	)
	if err != nil {
		server.HandlerErr(ctx, errCode, err)
		return
	}

	taskPayload := &common.PayloadSendVerifyEmail{
		Email:     payload.Email,
		FirstName: payload.FirstName,
		LastName:  payload.LastName,
		LangCode:  "ru",
		JWTToken:  token,
	}
	opts := []asynq.Option{
		asynq.MaxRetry(10),
		asynq.ProcessIn(10 * time.Second),
		asynq.Queue(scheduler.QueueCritical),
	}
	err = handler.taskDistributor.DistributeTaskSendVerifyEmail(ctx, taskPayload, opts...)

	if err != nil {
		log.Info().Err(err).Msg(fmt.Sprintf("distribute task send verify email err: %v", err))
	}
	ctx.JSON(http.StatusOK, server.Response(nil, errCode, nil))
}

func (handler *AuthHandler) EmailConfirmation(ctx *gin.Context) {
	token := ctx.Query("token")
	if token == "" {
		ctx.JSON(http.StatusBadRequest, server.Response(nil, server.INVALID_URL_PARAM_ERR_CODE, nil))
		return
	}
	errCode, err := handler.usecase.EmailConfirmation(ctx, token)
	if err != nil {
		server.HandlerErr(ctx, errCode, err)
		return
	}
	ctx.JSON(http.StatusOK, server.Response(nil, 0, nil))
}

func (handler *AuthHandler) SignInUser(ctx *gin.Context) {
	var payload *entities.SignInReq

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, server.Response(err, server.INVALID_DATA_ERR_CODE, nil))
		return
	}
	user, groups, errCode, err := handler.usecase.GetUser(ctx, payload)
	if err != nil {
		server.AuthHandlerErr(ctx, errCode, err)
		return
	}
	accessToken, _, errCode, err := handler.usecase.CreateAccessAndRefreshToken(user, groups, "access")
	if err != nil {
		server.HandlerErr(ctx, errCode, err)
		return
	}
	refreshToken, _, errCode, err := handler.usecase.CreateAccessAndRefreshToken(user, groups, "refresh")
	if err != nil {
		server.HandlerErr(ctx, errCode, err)
		return
	}
	response := gin.H{"access_token": accessToken, "refresh_token": refreshToken}
	ctx.JSON(http.StatusOK, server.Response(nil, errCode, response))
}

func (handler *AuthHandler) RefreshAccessToken(ctx *gin.Context) {
	var payload *entities.RefreshTokenReq
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, server.Response(err, server.INVALID_DATA_ERR_CODE, nil))
		return
	}

	accessToken, errCode, err := handler.usecase.RefreshAccessToken(ctx, payload.RefreshToken)
	if err != nil {
		server.HandlerErr(ctx, errCode, err)
		return
	}
	response := gin.H{"access_token": accessToken}
	ctx.JSON(http.StatusOK, server.Response(nil, errCode, response))
}

func (handler *AuthHandler) ChangePassword(ctx *gin.Context) {
	var payload *entities.ChangePasswordReq
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, server.Response(err, server.INVALID_DATA_ERR_CODE, nil))
		return
	}
	jwtPayload, exists := jwt_token.GetJWTPayload(ctx)
	if !exists {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, server.Response(nil, server.UNKNOWN_ERROR_CODE, nil))
		return
	}
	errCode, err := handler.usecase.ChangePassword(ctx, payload, jwtPayload.UserId)
	if err != nil {
		server.HandlerErr(ctx, errCode, err)
		return
	}
	ctx.JSON(http.StatusOK, server.Response(nil, errCode, nil))
}
