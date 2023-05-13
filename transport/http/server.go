package http

import (
	"database/sql"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	httpLog "github.com/labstack/gommon/log"
	domain "github.com/novan/golang-api-server/domain/errors"
	mw "github.com/novan/golang-api-server/transport/http/middleware"
	"github.com/novan/golang-api-server/transport/http/model"
	httpValidator "github.com/novan/golang-api-server/transport/http/validator"
	"github.com/novan/golang-api-server/util"
)

func Run() *echo.Echo {
	e := echo.New()
	// Set Bundle MiddleWare
	e.Use(middleware.RequestID())
	e.Pre(middleware.RemoveTrailingSlash())
	e.Use(middleware.RecoverWithConfig(middleware.RecoverConfig{
		StackSize:    1 << 10,
		LogLevel:     httpLog.ERROR,
		LogErrorFunc: HttpLogError,
	}))
	// e.Use(middleware.Gzip())
	e.Use(middleware.GzipWithConfig(middleware.GzipConfig{
		Skipper: func(c echo.Context) bool {
			return strings.Contains(c.Request().URL.Path, "swagger")
		},
	}))
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:  []string{"*"},
		AllowHeaders:  []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization, echo.HeaderContentLength, echo.HeaderAcceptEncoding, echo.HeaderAccessControlAllowOrigin, echo.HeaderAccessControlAllowHeaders, echo.HeaderContentDisposition, "X-Request-Id", "device-id", "X-Summary", "X-Account-Number", "X-Business-Name", "client-secret", "X-CSRF-Token"},
		ExposeHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization, echo.HeaderContentLength, echo.HeaderAcceptEncoding, echo.HeaderAccessControlAllowOrigin, echo.HeaderAccessControlAllowHeaders, echo.HeaderContentDisposition, "X-Request-Id", "device-id", "X-Summary", "X-Account-Number", "X-Business-Name", "client-secret", "X-CSRF-Token"},
		AllowMethods:  []string{echo.GET, echo.HEAD, echo.PUT, echo.PATCH, echo.POST, echo.DELETE},
	}))
	e.Use(mw.AccessLog)
	// e.Use(middleware.BodyDump(DebugHTTPBody))

	echo.NotFoundHandler = NotFoundErrorHandler
	e.HTTPErrorHandler = CustomErrorHandler
	e.Validator = httpValidator.NewCustomValidator()

	return e

}

func CustomErrorHandler(err error, c echo.Context) {

	var httpError int
	appErr, appErrOk := err.(*domain.AppError)
	he, heOk := err.(*echo.HTTPError)

	if err == sql.ErrNoRows {
		appErr = domain.NewAppHttpError(err, http.StatusNotFound)
		httpError = http.StatusNotFound
	} else if appErrOk {
		switch appErr.Type {
		case domain.BadRequest:
			httpError = http.StatusBadRequest
		case domain.NotFound:
			httpError = http.StatusNotFound
		case domain.ValidationError:
			httpError = http.StatusBadRequest
		case domain.ResourceAlreadyExists:
			httpError = http.StatusConflict
		case domain.NotAuthenticated:
			httpError = http.StatusUnauthorized
		case domain.NotAuthorized:
			httpError = http.StatusForbidden
		case domain.Forbidden:
			httpError = http.StatusForbidden
		default:
			httpError = http.StatusInternalServerError
		}
	} else if heOk && he.Internal == nil {
		appErr = domain.NewAppHttpError(err, he.Code)
		httpError = he.Code
	} else {
		appErr = domain.NewAppError(err, domain.InternalError)
		httpError = http.StatusInternalServerError
	}

	resp := model.NewErrorResponse(httpError, appErr.Type, err)

	_ = resp.JSON(c)
}

func NotFoundErrorHandler(c echo.Context) error {
	appErr := domain.NewAppErrorWithType(domain.NotFound)
	// render your 404 page
	resp := model.NewErrorResponse(http.StatusNotFound, appErr.Type, appErr.Err)
	return c.JSON(http.StatusNotFound, resp)
}

func HttpLogError(c echo.Context, err error, stack []byte) error {

	ctx := c.Request().Context()
	util.Log.WithContext(ctx).WithError(err).Fatalf("Recovered from Exception >> %s | %s", err.Error(), stack)

	return err
}

func DebugHTTPBody(c echo.Context, reqBody, resBody []byte) {
	ctx := c.Request().Context()
	util.Log.WithContext(ctx).Debugf("DebugHTTPBody | Request: %s | Response: %s", reqBody, resBody)
}
