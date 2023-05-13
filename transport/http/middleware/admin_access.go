package middleware

import (
	"net/http"

	"github.com/labstack/echo/v4"

	domainError "github.com/novan/golang-api-server/domain/errors"
	repo "github.com/novan/golang-api-server/repo/redis"
	"github.com/novan/golang-api-server/transport/http/model"
	"github.com/novan/golang-api-server/util"
)

func AdminAccess(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		var sess = c.Get(util.CONTEXT_SESSION)
		if sess != nil {
			user := sess.(*repo.User)
			if user.UserType == util.USERTYPE_ADMIN {
				return next(c)
			}
		}

		resp := model.Response{
			Code:    http.StatusUnauthorized,
			Status:  http.StatusText(http.StatusUnauthorized),
			Message: domainError.UnauthorizedAccessError().Error(),
		}
		return c.JSON(http.StatusUnauthorized, resp)
	}
}
