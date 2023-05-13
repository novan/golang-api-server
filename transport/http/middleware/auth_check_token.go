package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"

	"github.com/novan/golang-api-server/dic"
	domainError "github.com/novan/golang-api-server/domain/errors"
	"github.com/novan/golang-api-server/repo/redis"
	"github.com/novan/golang-api-server/transport/http/model"
	"github.com/novan/golang-api-server/util"
)

func AuthCheckToken(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		if ctx == nil {
			ctx = context.Background()
		}

		token := c.Request().Header.Get(echo.HeaderAuthorization)
		if token == "" {
			resp := model.Response{
				Code:    http.StatusUnauthorized,
				Status:  http.StatusText(http.StatusUnauthorized),
				Message: domainError.InvalidTokenError().Error(),
			}
			return c.JSON(http.StatusUnauthorized, resp)
		}
		token = strings.Replace(token, "Bearer ", "", -1)
		token = strings.Replace(token, "bearer ", "", -1)

		session := dic.Container.Get(dic.SessionRepository).(redis.SessionRepositoryInterface)
		val, err := session.GetUser(ctx, token)
		if err != nil {
			resp := model.Response{
				Code:    http.StatusUnauthorized,
				Status:  http.StatusText(http.StatusUnauthorized),
				Message: domainError.InvalidTokenError().Error(),
			}
			return c.JSON(http.StatusUnauthorized, resp)
		}

		c.Set(util.CONTEXT_TOKEN, token)
		c.Set(util.CONTEXT_SESSION, val)
		return next(c)
	}
}
