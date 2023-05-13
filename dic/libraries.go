package dic

import (
	"github.com/sarulabs/di/v2"

	// Repositories
	dbMysql "github.com/novan/golang-api-server/repo/mysql"
	dbRedis "github.com/novan/golang-api-server/repo/redis"

)

func RegisterLibraries(builder *di.Builder) {

	builder.Add(di.Def{
		Name: Db,
		Build: func(ctn di.Container) (interface{}, error) {
			return dbMysql.Connect(), nil
		},
	})

	builder.Add(di.Def{
		Name: Redis,
		Build: func(ctn di.Container) (interface{}, error) {
			return dbRedis.OpenClient(), nil
		},
	})

}
