package util

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"os"
	"strings"

	// "time"

	"github.com/labstack/echo/v4"
)

func StatusText(code int) string {
	if code == http.StatusOK {
		return "Success"
	} else {
		return http.StatusText(code)
	}
}

func FormatRequest(r *http.Request) string {
	// Create return string
	var request []string
	// Add the request string
	url := fmt.Sprintf("%v %v %v", r.Method, r.URL, r.Proto)
	request = append(request, url)
	// Add the host
	request = append(request, fmt.Sprintf("Host: %v", r.Host))
	// Loop through headers
	for name, headers := range r.Header {
		name = strings.ToLower(name)
		for _, h := range headers {
			request = append(request, fmt.Sprintf("%v: %v", name, h))
		}
	}

	// If this is a POST, add post data
	if r.Method == "POST" {
		r.ParseForm()
		request = append(request, "\n")
		request = append(request, r.Form.Encode())
	}
	// Return the request as a string
	return strings.Join(request, "\n")
}

func FindRoute(path string, method string, c echo.Context) string {
	for _, route := range c.Echo().Routes() {
		if route.Path == path && route.Method == method {
			return route.Name
		}
	}
	return ""
}

func HttpRequestDebug(r *http.Request) {
	if os.Getenv("APP_ENV") != "production" {
		data, err := httputil.DumpRequest(r, true)
		if err == nil {
			Log.Debugf("Dump Request:\n%s\n\n", data)
		} else {
			Log.Debugf("Dump Request error:\n%s\n\n", err)
		}
	}
}

func HttpResponseDebug(r *http.Response) {
	if os.Getenv("APP_ENV") != "production" {
		data, err := httputil.DumpResponse(r, true)
		if err == nil {
			Log.Debugf("Dump Response:\n%s\n\n", data)
		} else {
			Log.Debugf("Dump Response error:\n%s\n\n", err)
		}
	}
}
