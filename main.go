package main

import (
	"github.com/gsabadini/go-bank-transfer/infrastructure/database"
	"os"
	"time"

	"github.com/gsabadini/go-bank-transfer/infrastructure"
	"github.com/gsabadini/go-bank-transfer/infrastructure/logger"
	"github.com/gsabadini/go-bank-transfer/infrastructure/validator"
	"github.com/gsabadini/go-bank-transfer/infrastructure/web"
)

func main() {
	var app = infrastructure.NewConfig()

	app.Name(os.Getenv("APP_NAME"))

	app.Logger(logger.InstanceLogrusLogger)

	app.Validator(validator.InstanceGoPlayground)

	app.DbSQL(database.InstancePostgres).
		DbSQLCtxTimeout(5 * time.Second)

	app.DbNoSQL(database.InstanceMongoDB).
		DbNoSQLCtxTimeout(5 * time.Second)

	app.WebServerPort(os.Getenv("APP_PORT")).
		WebServer(web.InstanceGorillaMux).
		Start()
}
