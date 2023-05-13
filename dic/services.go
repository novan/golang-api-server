package dic

import (
	"github.com/sarulabs/di/v2"

	// Services
	"github.com/novan/golang-api-server/domain/service"

	// Repositories
	repo "github.com/novan/golang-api-server/repo/mysql"
	session "github.com/novan/golang-api-server/repo/redis"
)

func RegisterServices(builder *di.Builder) {

	builder.Add(di.Def{
		Name: AccountService,
		Build: func(ctn di.Container) (interface{}, error) {
			return service.NewAccountService(
				ctn.Get(UserRepository).(repo.UserRepositoryInterface),
				ctn.Get(TokenRepository).(repo.TokenRepositoryInterface),
				ctn.Get(SessionRepository).(session.SessionRepositoryInterface),
			), nil
		},
	})
}
