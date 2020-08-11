package main

import (
	"os"
	"time"

	"github.com/gsabadini/go-bank-transfer/infrastructure"
	"github.com/gsabadini/go-bank-transfer/infrastructure/database"
	"github.com/gsabadini/go-bank-transfer/infrastructure/logger"
	"github.com/gsabadini/go-bank-transfer/infrastructure/validator"
	"github.com/gsabadini/go-bank-transfer/infrastructure/web"
)

func main() {
	var app = infrastructure.NewConfig()

	app.Name(os.Getenv("APP_NAME"))

	app.ContextTimeout(10 * time.Second)

	app.Logger(logger.InstanceLogrusLogger)

	app.Validator(validator.InstanceGoPlayground)

	app.DbSQL(database.InstancePostgres)

	app.DbNoSQL(database.InstanceMongoDB)

	app.WebServerPort(os.Getenv("APP_PORT")).WebServer(web.InstanceGin).Start()
}
