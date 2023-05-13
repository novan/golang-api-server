package http

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"

	"github.com/novan/golang-api-server/dic"
	_ "github.com/novan/golang-api-server/docs"
	controller "github.com/novan/golang-api-server/transport/http/api/v1"
	"github.com/novan/golang-api-server/transport/http/middleware"
	"github.com/novan/golang-api-server/util"
)

// Route for mapping from json file
type Route struct {
	Name       string   `json:"name"`
	Path       string   `json:"path"`
	TargetHost string   `json:"target_host"`
	TargetPath string   `json:"target_path"`
	Method     string   `json:"method"`
	Module     string   `json:"module"`
	Endpoint   string   `json:"handler"`
	Middleware []string `json:"middleware"`
}

var (
	endpoint          map[string]echo.HandlerFunc
	middlewareHandler map[string]echo.MiddlewareFunc
)

type RoutesFactory struct {
	group *echo.Group
}

func NewRoutesFactory(e *echo.Echo) *RoutesFactory {

	e.GET("/health", func(c echo.Context) error {
		return c.NoContent(200)
	})

	if os.Getenv("APP_ENV") != "production" {
		e.GET("/swagger/*", echoSwagger.WrapHandler)
	}

	// Router binding
	apiV1 := e.Group("/api/v1")

	return &RoutesFactory{group: apiV1}
}

func (rf *RoutesFactory) Init() {

	rf.initEndpoint()
	rf.initMiddleware()

	routes := rf.LoadRoutes("./transport/http/gate/")
	for _, route := range routes {

		util.Log.Infof("HTTP Server | Adding route: %+v", route)
		if endpoint[route.Endpoint] == nil {
			util.Log.Fatalf("HTTP Server | Invalid endpoint: %s", route.Endpoint)
		}
		r := rf.group.Add(route.Method, route.Path, endpoint[route.Endpoint], rf.chainMiddleware(route)...)
		r.Name = route.Name
	}

}

func (rf *RoutesFactory) initEndpoint() {
	account := dic.Container.Get(dic.AccountController).(*controller.AccountController)

	endpoint = map[string]echo.HandlerFunc{
		// account
		"account.login":         account.Login,
		"account.signup":        account.Signup,
		"account.refresh":       account.Refresh,
		"account.toggle-active": account.ToggleActive,
	}
}

func (rf *RoutesFactory) initMiddleware() {
	middlewareHandler = map[string]echo.MiddlewareFunc{
		"auth":       middleware.AuthCheckToken,
		"admin_only": middleware.AdminAccess,
	}
}

func (rf *RoutesFactory) LoadRoutes(filePath string) []Route {
	var routes []Route
	files, err := ioutil.ReadDir(filePath)
	if err != nil {
		util.Log.Fatal("Failed to load file: %v", err)
	}
	for _, file := range files {
		fileName := file.Name()
		fileExt := filepath.Ext(filePath + "/" + fileName)
		if fileExt != ".json" {
			continue
		}
		byteFile, err := ioutil.ReadFile(filePath + "/" + file.Name())
		if err != nil {
			util.Log.Fatalf("Failed to load file: %+v", err)
		}
		var tmp []Route
		if err := util.Json.Unmarshal(byteFile, &tmp); err != nil {
			util.Log.Fatalf("Failed to marshal file %s: %v", file.Name(), err)
		}
		routes = append(routes, tmp...)
	}

	return routes
}

func (rf *RoutesFactory) chainMiddleware(route Route) []echo.MiddlewareFunc {
	var mwHandlers []echo.MiddlewareFunc

	for _, v := range route.Middleware {
		if middlewareHandler[v] == nil {
			util.Log.Fatalf("Invalid middleware: %s", v)
		}
		mwHandlers = append(mwHandlers, middlewareHandler[v])
	}
	return mwHandlers
}
