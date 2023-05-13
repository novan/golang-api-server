package redis

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/novan/golang-api-server/util"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type SessionRepositoryTestSuite struct {
	suite.Suite
	repo SessionRepositoryInterface
	ctx  context.Context
}

func (s *SessionRepositoryTestSuite) SetupTest() {
	fmt.Println("Initializing SessionRepositoryTestSuite")
	util.Env("../../../")
	fmt.Printf("APP_NAME: %s\n", os.Getenv("APP_NAME"))
	s.ctx = context.Background()
	s.repo = NewSessionRepository(OpenClient())
}

func (s *SessionRepositoryTestSuite) TestStoreUser() {
	now := time.Now()
	phone := gofakeit.Phone()
	token := gofakeit.UUID()
	user := User{
		ID:        gofakeit.Number(1, 100),
		UserType:  "PUBSCTR",
		Email:     gofakeit.Email(),
		Mobile:    &phone,
		LastLogin: &now,
		IsActive:  true,
	}
	fmt.Printf("Testing Store User | User: %+v\n", user)

	err := s.repo.StoreUser(s.ctx, token, &user)

	assert.Nil(s.T(), err)

	result, err := s.repo.GetUser(s.ctx, token)
	fmt.Printf("Testing Store User | Get User: %+v\n", result)

	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), result)
	assert.NotEmpty(s.T(), result.ID)
}

func TestSessionRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(SessionRepositoryTestSuite))
}
