package service

import (
	"context"
	"database/sql"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/jmoiron/sqlx/types"
	"github.com/novan/golang-api-server/domain/builder"
	"github.com/novan/golang-api-server/domain/entity"
	domainError "github.com/novan/golang-api-server/domain/errors"
	data "github.com/novan/golang-api-server/repo/mysql"
	"github.com/novan/golang-api-server/repo/mysql/schema"
	session "github.com/novan/golang-api-server/repo/redis"
	"github.com/novan/golang-api-server/util"
	"github.com/novan/golang-api-server/util/crypt"
	uuid "github.com/satori/go.uuid"
)

type AccountServiceInterface interface {
	Login(ctx context.Context, username string, password string) (*entity.UserToken, error)
	Signup(ctx context.Context, signup entity.SignupRequest, session *session.User) (*entity.User, error)
	RefreshToken(ctx context.Context, session *session.User, token string, refreshToken string) (*entity.UserToken, error)
	ToggleActive(ctx context.Context, userID int, isActive bool, user *session.User) error
}

type AccountService struct {
	repoUser  data.UserRepositoryInterface
	repoToken data.TokenRepositoryInterface
	session   session.SessionRepositoryInterface
	builder   *builder.AccountBuilder
}

func NewAccountService(
	repoUser data.UserRepositoryInterface,
	repoToken data.TokenRepositoryInterface,
	s session.SessionRepositoryInterface,
) *AccountService {
	return &AccountService{
		repoUser:  repoUser,
		repoToken: repoToken,
		session:   s,
		builder:   builder.NewAccountBuilder(),
	}
}

func (s *AccountService) Login(ctx context.Context, email string, password string) (*entity.UserToken, error) {

	user, err := s.repoUser.FindByEmail(ctx, strings.TrimSpace(email))
	if err != nil {
		return nil, domainError.NewAppError(err, domainError.RepositoryError)
	}

	if !crypt.Verify(user.Password, password) {
		return nil, domainError.InvalidUserLoginError()
	}

	genToken, genRefreshToken, err := s.generateToken(ctx, user)
	if err != nil {
		return nil, err
	}

	err = s.setupSession(ctx, strings.TrimSpace(email), genToken, user)
	if err != nil {
		return nil, err
	}

	resp := entity.UserToken{
		Token:        genToken,
		RefreshToken: genRefreshToken,
	}

	return &resp, nil
}

func (s *AccountService) generateToken(ctx context.Context, user *schema.User) (genToken string, genRefreshToken string, err error) {
	// creating access token
	genRefreshToken = uuid.NewV4().String()
	signingKey := []byte(os.Getenv("JWT_SECRET"))
	expire, _ := strconv.Atoi(os.Getenv("JWT_LIFETIME"))

	tokenID := uuid.NewV4()

	claims := entity.CustomClaims{
		Email:    user.Email,
		UserID:   user.ID,
		UserType: user.UserType.ToString(),
		StandardClaims: jwt.StandardClaims{
			Id:        tokenID.String(),
			NotBefore: time.Now().Unix(),
			ExpiresAt: time.Now().Add(time.Duration(expire) * time.Second).Unix(),
			Issuer:    os.Getenv("JWT_ISSUER"),
			Audience:  os.Getenv("JWT_AUDIENCE"),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedString, err := token.SignedString(signingKey)
	if err != nil {
		return "", "", domainError.NewAppError(err, domainError.TokenGeneratorError)
	}
	tmpToken, err := s.repoToken.GetToken(ctx, genRefreshToken)
	if err == sql.ErrNoRows {
		err = s.repoToken.CreateToken(ctx, genRefreshToken, tokenID.String(), user.ID)
	} else {
		tmpToken.JwtID = tokenID.String()
		err = s.repoToken.Update(ctx, *tmpToken)
	}
	if err != nil {
		return "", "", domainError.NewAppError(err, domainError.RepositoryError)
	}
	return signedString, genRefreshToken, nil
}

func (s *AccountService) Signup(ctx context.Context, signup entity.SignupRequest, session *session.User) (*entity.User, error) {

	if session == nil {
		signup.UserType = util.USERTYPE_USER
	} else if session.UserType != util.USERTYPE_ADMIN {
		return nil, domainError.UnauthorizedAccessError()
	}

	// check for existing user's email
	user, err := s.repoUser.FindByEmail(ctx, signup.Email)
	if err != nil && err != sql.ErrNoRows {
		util.Log.WithContext(ctx).WithError(err).Errorf("SignUp | Failed creating a person data | Error: %s", err.Error())
		return nil, err
	}
	if user != nil {
		return nil, domainError.EmailAlreadyExistError()
	}

	// password hasing
	hashed, err := crypt.Hash(signup.Password)
	if err != nil {
		util.Log.WithContext(ctx).WithError(err).Errorf("SignUp | Failed creating a password hash | Error: %s", err.Error())
		return nil, domainError.GeneratePasswordError()
	}

	// create a new user
	userModel := schema.User{
		UserType: signup.UserType,
		Email:    signup.Email,
		Password: hashed,
		Mobile:   signup.Mobile,
		IsActive: types.BitBool(true),
	}

	now := time.Now()
	userModel.CreatedAt = now
	userModel.UpdatedAt = now

	if session != nil {
		userModel.CreatedBy = session.ID
		userModel.UpdatedBy = session.ID
	} else {
		userModel.CreatedBy = 0
		userModel.UpdatedBy = 0
	}

	userID, err := s.repoUser.Create(ctx, userModel)
	if err != nil {
		util.Log.WithContext(ctx).WithError(err).Errorf("SignUp | Failed creating a user data | Error: %s", err.Error())
		return nil, err
	}
	userModel.ID = userID
	if userModel.CreatedBy == 0 {
		userModel.CreatedBy = userID
		userModel.UpdatedBy = userID
		s.repoUser.Update(ctx, userModel)
	}

	return s.builder.ToDomainModel(&userModel), nil
}

func (s *AccountService) RefreshToken(ctx context.Context, session *session.User, oldToken string, refreshToken string) (*entity.UserToken, error) {
	token, err := s.repoToken.GetToken(ctx, refreshToken)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, domainError.InvalidRefreshTokenError()
		}
		return nil, err
	}

	// validating the token
	if token.UserID != session.ID {
		return nil, domainError.NotMatchRefreshTokenError()
	}
	now := time.Now()
	if token.ExpiredAt.Before(now) {
		return nil, domainError.ExpiredRefreshTokenError()
	}
	if token.Invalidated {
		return nil, domainError.InvalidRefreshTokenError()
	}
	if token.IsUsed {
		return nil, domainError.UsedRefreshTokenError()
	} else {
		token.IsUsed = true
		token.Invalidated = true
		err := s.repoToken.Update(ctx, *token)
		if err != nil {
			util.Log.WithContext(ctx).WithError(err).Errorf("RefreshToken | Failed updating the token data | Error: %s", err.Error())
			return nil, err
		}
	}
	user, err := s.repoUser.FindByID(ctx, token.UserID)
	if err != nil {
		return nil, err
	}

	// generating new token
	newToken, newRefreshToken, _ := s.generateToken(ctx, user)
	sessionUser := s.builder.FromRepoToSessionModel(user)

	// store current session with new token and remove the old one
	err = s.session.StoreUser(ctx, newToken, sessionUser)
	if err != nil {
		util.Log.WithContext(ctx).WithError(err).Errorf("RefreshToken | Failed storing user session | Error: %s", err.Error())
	}
	_ = s.session.RemoveUser(ctx, oldToken)

	return &entity.UserToken{
		Token:        newToken,
		RefreshToken: newRefreshToken,
	}, nil
}

func (s *AccountService) setupSession(ctx context.Context, email string, genToken string, user *schema.User) (err error) {
	now := time.Now()
	user.LastLogin = &now
	_ = s.repoUser.Update(ctx, *user)
	if err != nil {
		util.Log.WithContext(ctx).WithError(err).Warnf("SetupSession | Failed updating user ID: %d | Error: %s", user.ID, err.Error())
	}

	sessionUser := s.builder.FromRepoToSessionModel(user)

	err = s.session.StoreUser(ctx, genToken, sessionUser)
	if err != nil {
		util.Log.WithContext(ctx).WithError(err).Errorf("Login | Failed storing user session | Error: %s", err.Error())
		return
	}

	return
}

func (s *AccountService) ToggleActive(ctx context.Context, userID int, isActive bool, user *session.User) error {
	u, err := s.repoUser.FindByID(ctx, userID)
	if err != nil {
		return err
	}
	u.IsActive = types.BitBool(isActive)
	u.UpdatedAt = time.Now()
	u.UpdatedBy = user.ID
	return s.repoUser.Update(ctx, *u)
}
