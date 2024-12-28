package middleware

import (
	"github.com/gin-gonic/gin"
	db "job_search_platform/internal/gateway_mrc/db/sqlc"
	"job_search_platform/pkg/helpers/server"
	"net/http"
)

func SessionMiddleware(store db.Store) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		uuidSession, err := ctx.Cookie("session_id")
		if err != nil {
			uuidSession = ""
		}
		userAgent := ctx.Request.UserAgent()
		clientIp := ctx.ClientIP()
		sessionArgs := db.RequestArgs{
			IpAddress: clientIp,
			UserAgent: userAgent,
		}
		session, created, errCode, err := store.GetOrCreateClientSession(ctx, uuidSession, sessionArgs)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, server.Response(err, errCode, nil))
			return
		}
		if !created && session.IsBlocked.Bool {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, server.Response(err, server.SESSION_BLOCKED_ERR_CODE, nil))
			return
		}
		if !created {
			_, err = store.UpdateSessionLastActive(ctx, uuidSession)
		}
		SessionId, _ := session.ID.Value()
		if created {
			ctx.SetCookie("session_id", SessionId.(string), 60*60*24, "/", "localhost", false, false)
		}
		//ctx.Set()
		ctx.Set("sessionId", SessionId.(string))
		ctx.Next()
	}
}
