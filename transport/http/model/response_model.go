package model

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	domainError "github.com/novan/golang-api-server/domain/errors"
	httpValidator "github.com/novan/golang-api-server/transport/http/validator"
	"gopkg.in/go-playground/validator.v9"
)

// Response struct
type Response struct {
	Code       int                    `json:"statusCode"`
	Status     string                 `json:"status"`
	Message    string                 `json:"message,omitempty"`
	Payload    interface{}            `json:"payload"`
	Errors     interface{}            `json:"errors,omitempty"`
	Pagination *Pagination            `json:"pagination,omitempty"`
	Header     map[string]interface{} `json:"-"`
	HttpCode   int                    `json:"-"`
}

func (r *Response) ToString() string {
	resp, _ := json.Marshal(r)
	return string(resp)
}

// WithPagination set response with pagination
func (r *Response) WithPagination(c echo.Context, pagination Pagination) *Response {
	r.Pagination = &pagination
	page := r.Pagination.CurrentPage
	u := c.Request().URL
	if r.Pagination.HasNextPage() {
		q := u.Query()
		q.Set("page", fmt.Sprintf("%d", page+1))
		u.RawQuery = q.Encode()
		r.Pagination.Next = u.String()
	}
	if page > 1 {
		q := u.Query()
		q.Set("page", fmt.Sprintf("%d", page-1))
		u.RawQuery = q.Encode()
		r.Pagination.Prev = u.String()
	}
	return r
}

// JSON render response as JSON
func (r *Response) JSON(c echo.Context) error {
	for k, v := range r.Header {
		c.Response().Header().Set(k, fmt.Sprintf("%s,%v", c.Response().Header().Get(k), v))
	}
	if r.HttpCode > 0 {
		return c.JSON(int(r.HttpCode), r)
	}
	return c.JSON(r.Code, r)
}

func (r *Response) SetMessage(msg string) {
	r.Message = msg
}

func NewSuccessResponseWithoutData() *Response {
	return &Response{
		HttpCode: http.StatusOK,
		Code:     http.StatusOK,
		Status:   "success",
		Message:  "Success",
	}
}

func NewSuccessResponse(data interface{}) *Response {
	return &Response{
		HttpCode: http.StatusOK,
		Code:     http.StatusOK,
		Status:   "success",
		Message:  "Success",
		Payload:  data,
	}
}

func NewErrorResponse(httpCode int, errType string, err error) *Response {
	var ve validator.ValidationErrors
	var ae httpValidator.ApiErrorResponse

	if errors.As(err, &ve) {
		return &Response{
			HttpCode: http.StatusBadRequest,
			Code:     http.StatusBadRequest,
			Status:   domainError.BadRequest,
			Message:  domainError.NewAppErrorWithType(domainError.BadRequest).Error(),
			Errors:   httpValidator.ValidatorMessageFormat(err),
		}
	} else if errors.As(err, &ae) {
		return &Response{
			HttpCode: ae.StatusCode,
			Code:     ae.StatusCode,
			Status:   ae.Status,
			Message:  ae.Message,
			Errors:   ae.Errors,
		}

	}
	return &Response{
		HttpCode: httpCode,
		Code:     httpCode,
		Status:   errType,
		Message:  err.Error(),
	}

}
