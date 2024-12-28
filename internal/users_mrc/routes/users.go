package routes

import (
	"github.com/gin-gonic/gin"
	"job_search_platform/internal/users_mrc/handlers"
)

type UsersRouter struct {
	handler handlers.UsersHandler
}

func NewUsersRouter(handler handlers.UsersHandler) *UsersRouter {
	return &UsersRouter{handler: handler}
}

func (r *UsersRouter) InitUsersRouter(public *gin.RouterGroup, private *gin.RouterGroup) {
	private.POST("/update", r.handler.UpdateUser)
}
