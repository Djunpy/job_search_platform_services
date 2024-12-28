package middleware

import (
	"github.com/gin-gonic/gin"
	"job_search_platform/pkg/helpers/server"
	"job_search_platform/pkg/jwt_token"
	"net/http"
	"strings"
)

func JWTDeserializer(tokenMaker jwt_token.Maker) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var accessToken string
		authorizationHeader := ctx.GetHeader("Authorization")
		fields := strings.Fields(authorizationHeader)

		if len(fields) != 2 || fields[0] != "Bearer" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, server.Response(nil, server.AUTH_HEADER_ERR_CODE, nil))
			return
		}
		accessToken = fields[1]
		jwtPayload, err := tokenMaker.VerifyToken(accessToken)
		if err != nil {
			statusCode := tokenMaker.GetErrorCode(err)
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, server.Response(err, statusCode, nil))
			return
		}

		ctx.Set("jwtTokenPayload", jwtPayload)
		ctx.Next()
	}
}
