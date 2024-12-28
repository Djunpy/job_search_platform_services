package handlers

import (
	"github.com/gin-gonic/gin"
	"job_search_platform/internal/users_mrc/entities"
	"job_search_platform/internal/users_mrc/usecases"
	"job_search_platform/pkg/helpers/server"
	"job_search_platform/pkg/jwt_token"
	"net/http"
)

type UsersHandler struct {
	usecase usecases.UsersUsecase
}

func NewUsersHandler(usecase usecases.UsersUsecase) UsersHandler {
	return UsersHandler{usecase: usecase}
}

func (handler *UsersHandler) UpdateUser(ctx *gin.Context) {
	var payload *entities.UserUpdate
	var err error
	if err = ctx.ShouldBindJSON(&payload); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, server.Response(err, server.INVALID_DATA_ERR_CODE, nil))
		return
	}
	jwtPayload, exists := jwt_token.GetJWTPayload(ctx)
	if !exists {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, server.Response(nil, server.UNKNOWN_ERROR_CODE, nil))
		return
	}
	errCode, err := handler.usecase.UpdateUser(ctx, *payload, jwtPayload.UserId)
	if err != nil {
		server.HandlerErr(ctx, errCode, err)
		return
	}
	ctx.JSON(http.StatusOK, server.Response(nil, errCode, nil))
}
