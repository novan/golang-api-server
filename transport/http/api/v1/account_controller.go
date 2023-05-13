package v1

import (
	"github.com/labstack/echo/v4"
	"github.com/novan/golang-api-server/domain/builder"
	_ "github.com/novan/golang-api-server/domain/entity"
	domain "github.com/novan/golang-api-server/domain/service"
	repo "github.com/novan/golang-api-server/repo/redis"
	"github.com/novan/golang-api-server/transport/http/model"
	"github.com/novan/golang-api-server/util"
)

type AccountController struct {
	service domain.AccountServiceInterface
	builder *builder.AccountBuilder
}

func NewAccountController(service domain.AccountServiceInterface) *AccountController {
	return &AccountController{
		service: service,
		builder: builder.NewAccountBuilder(),
	}
}

// Signup godoc
// @Summary Signup
// @Description Register a user
// @Tags Account
// @Accept application/json
// @Produce json
// @Param SignupRequest body model.SignupRequest true "Sign Up Request"
// @Router /account/signup [POST]
// @security Bearer
// @response 200 {object} model.Response{} "Success"
func (c *AccountController) Signup(e echo.Context) error {
	ctx := e.Request().Context()
	s := new(model.SignupRequest)
	if err := e.Bind(s); err != nil {
		return err
	}
	if err := e.Validate(s); err != nil {
		return err
	}

	// reject if user trying to create non-public user but not admin
	var user *repo.User
	session := e.Get(util.CONTEXT_SESSION)
	if session != nil {
		user = session.(*repo.User)
	}

	newUser, err := c.service.Signup(ctx, c.builder.FromSignupRequestToDomain(s), user)
	if err == nil {
		res := model.NewSuccessResponse(newUser)
		return res.JSON(e)
	}

	return err
}

// Login godoc
// @Summary Login
// @Description Authenticate a user
// @Tags Account
// @Accept application/json
// @Produce  json
// @Param LoginRequest body model.LoginRequest true "Login Request"
// @Router /account/login [POST]
// @security Bearer
// @response 200 {object} model.Response{Payload=entity.UserToken} "Success"
func (c *AccountController) Login(e echo.Context) error {
	ctx := e.Request().Context()
	s := new(model.LoginRequest)
	if err := e.Bind(s); err != nil {
		return err
	}
	if err := e.Validate(s); err != nil {
		return err
	}

	token, err := c.service.Login(ctx, s.Email, s.Password)
	if err == nil {
		res := model.NewSuccessResponse(token)
		return res.JSON(e)
	}

	return err
}

// Refresh godoc
// @Summary Refresh Token
// @Description Refresh an access token
// @Tags Account
// @Accept application/json
// @Produce  json
// @Param model.RefreshTokenRequest body model.RefreshTokenRequest true "Refresh token request"
// @Router /account/refresh [POST]
// @security Bearer
// @response 200 {object} model.Response{Payload=entity.UserToken} "Success"
func (c *AccountController) Refresh(e echo.Context) error {
	ctx := e.Request().Context()
	s := new(model.RefreshTokenRequest)
	if err := e.Bind(s); err != nil {
		return err
	}
	if err := e.Validate(s); err != nil {
		return err
	}

	// reject if user trying to create non-public user but not admin
	var user *repo.User
	session := e.Get(util.CONTEXT_SESSION)
	if session != nil {
		user = session.(*repo.User)
	}

	userToken, err := c.service.RefreshToken(ctx, user, e.Get(util.CONTEXT_TOKEN).(string), s.RefreshToken)
	if err == nil {
		res := model.NewSuccessResponse(userToken)
		return res.JSON(e)
	}

	return err
}

// ToggleActive godoc
// @Summary Toggling active/inactive user
// @Description Toggling active/inactive user
// @Tags Account
// @Accept application/json
// @Produce  json
// @Param ToggleActiveRequest body model.ToggleActiveRequest true "Request"
// @Router /account/toggle [PUT]
// @security Bearer
// @response 200 {object} model.Response{} "Success"
func (c *AccountController) ToggleActive(e echo.Context) error {
	ctx := e.Request().Context()
	s := new(model.ToggleActiveRequest)
	if err := e.Bind(s); err != nil {
		return err
	}
	if err := e.Validate(s); err != nil {
		return err
	}

	// reject if user trying to create non-public user but not admin
	var user *repo.User
	session := e.Get(util.CONTEXT_SESSION)
	if session != nil {
		user = session.(*repo.User)
	}

	err := c.service.ToggleActive(ctx, s.UserID, s.IsActive, user)
	if err == nil {
		res := model.NewSuccessResponseWithoutData()
		return res.JSON(e)
	}

	return err
}
