package dic

import (
	"github.com/sarulabs/di/v2"

	// Services
	"github.com/novan/golang-api-server/domain/service"

	// Controller
	controller "github.com/novan/golang-api-server/transport/http/api/v1"
)

func RegisterControllers(builder *di.Builder) {

	builder.Add(di.Def{
		Name: AccountController,
		Build: func(ctn di.Container) (interface{}, error) {
			return controller.NewAccountController(ctn.Get(AccountService).(service.AccountServiceInterface)), nil
		},
	})

}