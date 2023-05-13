package http

import (
	"context"
	"net/http"
	"net/url"

	"github.com/gomodule/oauth1/oauth"
)

type IHttp interface {
	FormatUrl(endpoint string) string
	RequestWithFormData(ctx context.Context, urlPath string, body map[string]string, headers map[string]string) (resp *http.Response, bodyBytes []byte, err error)
	RequestWithJSON(ctx context.Context, urlPath string, jsonData interface{}, headers map[string]string) (resp *http.Response, bodyBytes []byte, err error)
	SendForward(ctx context.Context, request *http.Request, targetPath string) (*http.Response, []byte, error)
	Send(ctx context.Context, request *http.Request) (*http.Response, []byte, error)
	SetOAuthCredentials(credentials oauth.Credentials)
	SetBearerToken(token string)
	SetAuthMode(authMode AuthorizationType)
	SetAuthorizationHeader(request *http.Request) error
	SignForm(request *http.Request, key string, secret string) url.Values
	CreateRequest(ctx context.Context, method string, urlPath string, params map[string]string, headers map[string]string) (request *http.Request, err error)
	CreateRequestJSON(ctx context.Context, urlPath string, jsonData interface{}, headers map[string]string) (request *http.Request, err error)
}
