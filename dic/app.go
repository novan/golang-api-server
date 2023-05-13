package dic

import (
	"github.com/sarulabs/di/v2"
)

var Builder *di.Builder
var Container di.Container

const Db = "db"
const Redis = "redis"

const AccountService = "service.account"

const SessionRepository = "repository.session"
const UserRepository = "repository.user"
const TokenRepository = "repository.token"

const AccountController = "controller.account"

func InitContainer() di.Container {
	builder := InitBuilder()
	Container = builder.Build()
	return Container
}

func InitBuilder() *di.Builder {
	Builder, _ = di.NewBuilder()
	RegisterLibraries(Builder)
	RegisterServices(Builder)
	RegisterRepositories(Builder)
	RegisterControllers(Builder)
	return Builder
}
