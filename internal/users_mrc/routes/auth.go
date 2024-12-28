package routes

import (
	"github.com/gin-gonic/gin"
	"job_search_platform/internal/users_mrc/handlers"
)

type AuthRouter struct {
	handler handlers.AuthHandler
}

func NewAuthRouter(handler handlers.AuthHandler) *AuthRouter {
	return &AuthRouter{handler: handler}
}

func (r *AuthRouter) InitAuthRouter(public *gin.RouterGroup, private *gin.RouterGroup) {
	public.POST("/sign-up", r.handler.SignUpUser)
	public.GET("/email-confirmation", r.handler.EmailConfirmation)
	public.POST("/sign-in", r.handler.SignInUser)
	public.POST("/refresh-token", r.handler.RefreshAccessToken)
	private.POST("/change-password", r.handler.ChangePassword)
}
