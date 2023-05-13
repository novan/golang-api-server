package dic

import (
	"github.com/go-redis/redis/v8"
	"github.com/jmoiron/sqlx"
	"github.com/sarulabs/di/v2"

	// Repositories
	repo "github.com/novan/golang-api-server/repo/mysql"
	session "github.com/novan/golang-api-server/repo/redis"
)

func RegisterRepositories(builder *di.Builder) {

	builder.Add(di.Def{
		Name: SessionRepository,
		Build: func(ctn di.Container) (interface{}, error) {
			return session.NewSessionRepository(ctn.Get(Redis).(*redis.Client)), nil
		},
	})

	builder.Add(di.Def{
		Name: UserRepository,
		Build: func(ctn di.Container) (interface{}, error) {
			return repo.NewUserRepository(ctn.Get(Db).(*sqlx.DB)), nil
		},
	})

	builder.Add(di.Def{
		Name: TokenRepository,
		Build: func(ctn di.Container) (interface{}, error) {
			return repo.NewTokenRepository(ctn.Get(Db).(*sqlx.DB)), nil
		},
	})

}
