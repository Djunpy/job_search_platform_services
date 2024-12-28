package server

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

var (
	// status codes
	SUCCESS_CODE                      int32 = 0  // Успешная операция
	INVALID_DATA_ERR_CODE             int32 = 27 // Неверные данные
	JWT_GENERATE_ERR_CODE             int32 = 2  // Ошибка генерации JWT
	GET_COOKIE_ERR_CODE               int32 = 3  // Ошибка получения Cookie
	TOKEN_VALIDATION_ERR_CODE         int32 = 4  // Ошибка валидации токена
	GENERATE_JWT_TOKEN_ERR_CODE       int32 = 5  // Ошибка генерации JWT токена
	MISSING_JWT_TOKEN_ERR_CODE        int32 = 6  // JWT токен отсутствует
	AUTH_HEADER_NOT_PROVIDED_ERR_CODE int32 = 7  // Заголовок авторизации не предоставлен
	AUTH_HEADER_FORMAT_ERR_CODE       int32 = 8  // Ошибка формата заголовка авторизации
	AUTH_HEADER_TYPE_ERR_CODE         int32 = 9  // Неверный тип заголовка авторизации
	AUTO_LOGOUT_ERR_CODE              int32 = 10 // Ошибка автоматического выхода
	JWT_EXPIRES_ERR_CODE              int32 = 11 // JWT токен истек
	USER_EXISTS_ERR_CODE              int32 = 12 // Пользователь уже существует
	INCORRECT_PASSWORD_ERR_CODE       int32 = 13 // Неверный пароль
	USER_NOT_EXISTS_ERR_CODE          int32 = 14 // Пользователь не существует
	EMPTY_FIELD_ERR_CODE              int32 = 15 // Поле пусто
	INVALID_URL_PARAM_ERR_CODE        int32 = 16 // Неверный параметр URL
	USER_INFO_NOT_FOUND_CODE          int32 = 17 // Информация о пользователе не найдена
	INVITE_CODE_USED_ERR_CODE         int32 = 18 // Пригласительный код уже использован
	SESSION_PARSING_ERR_CODE          int32 = 19 // Ошибка разбора сессии
	TOKEN_REFRESH_ERR_CODE            int32 = 20 // Ошибка обновления токена
	AUTH_HEADER_ERR_CODE              int32 = 21 // Общая ошибка заголовка авторизации
	SESSION_BLOCKED_ERR_CODE          int32 = 22 // Сессия пользователя заблокирована
	CREATING_REQUEST_ERR_CODE         int32 = 23
	PARSING_RESPONSE_ERR_CODE         int32 = 24
	SENDING_TOKEN_REFRESH_ERR_CODE    int32 = 25
	SESSION_NOT_FOUND_ERR_CODE        int32 = 26
	UNKNOWN_ERROR_CODE                int32 = 1
)

func Response(err error, code int32, body interface{}) gin.H {
	if err != nil {
		return gin.H{"error": err.Error(), "code": code, "body": body}
	}
	return gin.H{"error": nil, "code": code, "body": body}
}

func HandlerErr(ctx *gin.Context, errCode int32, err error) {
	if err != nil {
		ctx.JSON(http.StatusBadRequest, Response(err, errCode, nil))
	} else {
		ctx.JSON(http.StatusBadRequest, Response(nil, errCode, nil))
	}
}

func AuthHandlerErr(ctx *gin.Context, errCode int32, err error) {
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, Response(err, errCode, nil))
	} else {
		ctx.JSON(http.StatusBadRequest, Response(nil, errCode, nil))
	}
}
