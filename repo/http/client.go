package http

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	netHttp "net/http"
	"net/url"
	"strings"
	"time"

	"github.com/gomodule/oauth1/oauth"
	"github.com/labstack/echo/v4"
	domainError "github.com/novan/golang-api-server/domain/errors"
	"github.com/novan/golang-api-server/util"
)

type AuthorizationType string

const (
	AUTHORIZATION_BASIC  AuthorizationType = "Basic"
	AUTHORIZATION_OAUTH                    = "OAuth"
	AUTHORIZATION_BEARER                   = "Bearer"
)

type ApiClient struct {
	Client      netHttp.Client
	baseUrl     string
	authMode    AuthorizationType
	oauth       *oauth.Credentials
	bearerToken *string
}

func NewApiClient(baseUrl string) *ApiClient {
	tr := &netHttp.Transport{
		MaxIdleConns:       30,
		IdleConnTimeout:    60 * time.Second,
		DisableCompression: false,
	}
	return &ApiClient{
		Client:  netHttp.Client{Transport: tr},
		baseUrl: baseUrl,
		oauth:   nil,
	}
}

func (c *ApiClient) FormatUrl(endpoint string) string {
	return c.baseUrl + endpoint
}

func (c *ApiClient) SetOAuthCredentials(credentials oauth.Credentials) {
	c.oauth = &credentials
}

func (c *ApiClient) SetBearerToken(token string) {
	c.bearerToken = &token
}

func (c *ApiClient) SetAuthMode(authMode AuthorizationType) {
	c.authMode = authMode
}

func (c *ApiClient) SetAuthorizationHeader(request *netHttp.Request) (err error) {
	switch c.authMode {
	case AUTHORIZATION_OAUTH:
		client := oauth.Client{}
		client.Credentials = *c.oauth
		err = client.SetAuthorizationHeader(request.Header, c.oauth, request.Method, request.URL, request.Form)
	case AUTHORIZATION_BASIC:
		request.Header.Set("Authorization", fmt.Sprintf("Basic %s", base64.StdEncoding.EncodeToString([]byte(c.oauth.Token+":"+c.oauth.Secret))))
	case AUTHORIZATION_BEARER:
		if c.bearerToken != nil {
			request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", *c.bearerToken))
		} else {
			err = errors.New("Invalid bearer token")
		}
	}
	return err
}

func (c *ApiClient) SignForm(request *netHttp.Request, key string, secret string) url.Values {
	client := oauth.Client{}
	credentials := &oauth.Credentials{
		Token:  key,
		Secret: secret,
	}
	client.Credentials = *credentials
	params := request.URL.Query()
	client.SignForm(credentials, request.Method, request.URL.String(), params)
	return params
}

func (c *ApiClient) RequestWithFormData(ctx context.Context, urlPath string, body map[string]string, headers map[string]string) (resp *netHttp.Response, bodyBytes []byte, err error) {
	req, err := c.CreateRequest(ctx, netHttp.MethodPost, urlPath, body, headers)
	if err != nil {
		util.Log.WithContext(ctx).WithError(err).Errorf("ApiClient | Error creating a request | Error: %s", err.Error())
	}
	_, bodyBytes, err = c.Send(ctx, req)
	return
}

func (c *ApiClient) RequestWithJSON(ctx context.Context, urlPath string, jsonData interface{}, headers map[string]string) (resp *netHttp.Response, bodyBytes []byte, err error) {
	req, _ := c.CreateRequestJSON(ctx, urlPath, jsonData, headers)
	return c.Send(ctx, req)
}

func (c *ApiClient) SendForward(ctx context.Context, request *netHttp.Request, targetPath string) (*netHttp.Response, []byte, error) {
	var body io.Reader = nil
	data := request.Form

	multipartBody := &bytes.Buffer{}
	writer := multipart.NewWriter(multipartBody)

	if request.Method != netHttp.MethodGet {
		if len(request.Form) > 0 {
			body = strings.NewReader(data.Encode())
		} else if request.Body != nil {
			body = request.Body
		}

		if request.MultipartForm != nil {
			fileName := "file.png"
			for k, v := range data {
				if k != "file" {
					_ = writer.WriteField(k, v[0])
				}
				if k == "file_name" {
					fileName = v[0]
				}
			}
			fileHeader := request.MultipartForm.File["file"][0]
			file, _ := fileHeader.Open()
			part, _ := writer.CreateFormFile("file", fileName)
			io.Copy(part, file)
			body = multipartBody
		}
	}
	writer.Close()

	req, err := netHttp.NewRequestWithContext(ctx, request.Method, c.FormatUrl(targetPath), body)
	if err != nil {
		return nil, nil, err
	}

	if request.Method != netHttp.MethodGet {
		req.PostForm = data
	}

	if request.URL.Query() != nil {
		req.URL.RawQuery = request.URL.Query().Encode()
	}

	req.Header = request.Header.Clone()
	req.Header.Set(echo.HeaderXForwardedFor, request.RemoteAddr)
	if request.MultipartForm != nil {
		req.Header.Set("Content-Type", writer.FormDataContentType())
	}
	// log.Println("[REQUEST] detail request to forward to other http request")
	// log.Println(req)
	// bytes, _ := json.Marshal(&req)
	// util.Log.WithError(err).Errorf("[REQUEST] Detail request to forward to other http request: %s", string(bytes))
	return c.Send(ctx, req)
}

func (c *ApiClient) CreateRequest(ctx context.Context, method string, urlPath string, params map[string]string, headers map[string]string) (request *netHttp.Request, err error) {

	if method == netHttp.MethodGet {
		values := url.Values{}
		for k, v := range params {
			values.Set(k, v)
		}
		payload := values.Encode()
		if payload != "" {
			payload = "?" + payload
		}
		request, err = netHttp.NewRequestWithContext(ctx, method, c.FormatUrl(urlPath)+payload, nil)
	} else {
		payload := &bytes.Buffer{}
		var writer *multipart.Writer
		writer = multipart.NewWriter(payload)
		for k, v := range params {
			writer.WriteField(k, v)
		}
		err = writer.Close()

		request, err = netHttp.NewRequestWithContext(ctx, method, c.FormatUrl(urlPath), payload)
		request.Header.Set("Content-Type", writer.FormDataContentType())
	}
	if err != nil {
		util.Log.WithContext(ctx).WithError(err).Errorf("ApiClient | Failed creating request | Error: %s", err.Error())
		return nil, err
	}

	return request, nil
}

func (c *ApiClient) CreateRequestJSON(ctx context.Context, urlPath string, jsonData interface{}, headers map[string]string) (request *netHttp.Request, err error) {

	payload, err := json.Marshal(jsonData)
	if err != nil {
		return
	}

	request, err = netHttp.NewRequestWithContext(ctx, netHttp.MethodPost, c.FormatUrl(urlPath), bytes.NewBuffer(payload))
	if err == nil {
		for key, val := range headers {
			request.Header.Set(key, val)
		}
		request.Header.Set("Content-Type", "application/json")
	}

	return request, nil

}

func (c *ApiClient) Send(ctx context.Context, request *netHttp.Request) (*netHttp.Response, []byte, error) {

	request.ParseForm()

	if c.oauth != nil {
		err := c.SetAuthorizationHeader(request)
		if err != nil {
			util.Log.WithContext(ctx).WithError(err).Errorf("ApiClient | Failed set authorization header | Error: %s", err.Error())
		}
	}

	// util.HttpRequestDebug(request)
	util.Log.WithContext(ctx).Infof("ApiClient | HTTP Request: [%s] %s", request.Method, request.URL.String())
	resp, err := c.Client.Do(request)
	// util.HttpResponseDebug(resp)

	if err != nil {
		util.Log.WithContext(ctx).WithError(err).Errorf("ApiClient | Failed sending HTTP Request | Error: %s", err.Error())
		return nil, nil, domainError.InternalServerError(err.Error())
	}
	defer resp.Body.Close()

	var reader io.ReadCloser
	switch resp.Header.Get("Content-Encoding") {
	case "gzip":
		reader, err = gzip.NewReader(resp.Body)
		if err != nil {
			util.Log.WithContext(ctx).WithError(err).Errorf("ApiClient | Failed decompressing response body | Error: %s", err.Error())
			return nil, nil, domainError.InternalServerError(err.Error())
		}
		defer reader.Close()
	default:
		reader = resp.Body
	}

	body, err := ioutil.ReadAll(reader)
	if err != nil {
		util.Log.WithContext(ctx).WithError(err).Errorf("ApiClient | Failed reading response body | Error: %s", err.Error())
		return nil, nil, domainError.InternalServerError(err.Error())
	}
	util.Log.WithContext(ctx).Infof("ApiClient | HTTP Response: [%d] %s", resp.StatusCode, request.URL.String())

	return resp, body, nil
}
