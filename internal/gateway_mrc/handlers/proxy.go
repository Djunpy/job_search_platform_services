package handlers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/goccy/go-json"
	"io"
	"job_search_platform/internal/gateway_mrc/usecases"
	"job_search_platform/pkg/config"
	"job_search_platform/pkg/entities/common"
	"job_search_platform/pkg/helpers/server"
	"job_search_platform/pkg/jwt_token"
	"net/http"
)

type ProxyHandler struct {
	sessionsUsecase usecases.SessionsUsecase
	config          config.Config
	jwtMaker        jwt_token.Maker
}

func NewProxyHandler(jwtMaker jwt_token.Maker, sessionsUsecase usecases.SessionsUsecase, config config.Config) ProxyHandler {
	return ProxyHandler{sessionsUsecase: sessionsUsecase, jwtMaker: jwtMaker, config: config}
}

func (c *ProxyHandler) ProxySignInReq(ctx *gin.Context, target string) {
	var payload common.SignInResponse
	sessionIdStr, exists := ctx.Get("sessionId")
	if !exists {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, server.Response(nil, server.GET_COOKIE_ERR_CODE, nil))
		return
	}

	url := server.GetReqFullUrl(ctx, target)
	fmt.Print(ctx.Request.Body)
	resp, err := server.CreateAndSendRequest(
		ctx.Request.Method,
		url,
		ctx.Request.Body,
		ctx.Request.Header,
	)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, server.Response(err, server.CREATING_REQUEST_ERR_CODE, nil))
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read response"})
		return
	}

	if resp.StatusCode == http.StatusOK {
		err = json.Unmarshal(body, &payload)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, server.Response(err, server.PARSING_RESPONSE_ERR_CODE, nil))
			return
		}
		refreshToken := payload.Body.RefreshToken
		accessToken := payload.Body.AccessToken
		_, statusCode, err := c.sessionsUsecase.UpdateSession(ctx, sessionIdStr.(string), refreshToken, accessToken)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, server.Response(err, statusCode, nil))
			return
		}
	}
	if resp.StatusCode == http.StatusOK {
		ctx.JSON(resp.StatusCode, server.Response(nil, server.SUCCESS_CODE, nil))
	} else {
		ctx.Data(resp.StatusCode, resp.Header.Get("Content-Type"), body)
	}
}

func (c *ProxyHandler) ProxyCommonReq(ctx *gin.Context, target string) {
	sessionIdStr, exists := ctx.Get("sessionId")
	if !exists {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, server.Response(nil, server.GET_COOKIE_ERR_CODE, nil))
		return
	}

	session, errCode, err := c.sessionsUsecase.GetSession(ctx, sessionIdStr.(string))
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, server.Response(err, errCode, nil))
		return
	}
	url := server.GetReqFullUrl(ctx, target)
	headers := ctx.Request.Header.Clone()
	headers.Set("Authorization", fmt.Sprintf("Bearer %s", session.AccessToken.String))

	resp, err := server.CreateAndSendRequest(
		ctx.Request.Method,
		url,
		ctx.Request.Body,
		headers,
	)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, server.Response(err, server.CREATING_REQUEST_ERR_CODE, nil))
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, server.Response(err, server.CREATING_REQUEST_ERR_CODE, nil))
		return
	}

	ctx.Data(resp.StatusCode, resp.Header.Get("Content-Type"), body)
}

func (c *ProxyHandler) ProxyLogoutReq(ctx *gin.Context) {
	ctx.SetCookie("session_id", "", -1, "/", "localhost", false, true)
	ctx.Status(http.StatusOK)
}
