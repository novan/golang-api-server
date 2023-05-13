package http

import (
	"context"
	"fmt"
	"os"
	"testing"
	netHttp "net/http"

	// "github.com/brianvoe/gofakeit/v6"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"github.com/novan/golang-api-server/util"
)

type ClientHTTPTestSuite struct {
	suite.Suite
	ctx context.Context
	client IHttp
}

func (s *ClientHTTPTestSuite) SetupTest() {
	fmt.Println("Initializing ClientHTTPTestSuite")
	util.Env("../../")
	s.ctx = context.Background()
	s.client = NewApiClient(os.Getenv("WC_API_URL"))
}

func (s *ClientHTTPTestSuite) TestSignForm() {
	request, err := netHttp.NewRequest("POST", "/revo-admin/v1/generate_auth_cookie", nil)
	params := s.client.SignForm(request, os.Getenv("WC_CONSUMER_KEY"), os.Getenv("WC_CONSUMER_SECRET")) 
	
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), params.Get("oauth_consumer_key"), os.Getenv("WC_CONSUMER_KEY"))
	assert.Equal(s.T(), params.Get("oauth_signature_method"), "HMAC-SHA1")
}

func TestClientHTTPTestSuite(t *testing.T) {
	suite.Run(t, new(ClientHTTPTestSuite))
}
