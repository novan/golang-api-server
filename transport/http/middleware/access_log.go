package middleware

import (
	"context"

	"github.com/labstack/echo/v4"
	"github.com/novan/golang-api-server/util"
)

func AccessLog(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		req := c.Request()
		res := c.Response()

		ctx := req.Context()
		if ctx == nil {
			ctx = context.Background()
		}

		if err := next(c); err != nil {
			c.Error(err)
		}

		util.Log.WithContext(ctx).Debugf("HTTP Access | %s | [%s] %s (%d)", req.RemoteAddr, req.Method, req.RequestURI, res.Status)

		return nil
	}
}
