package main

import (
	"os"

	"github.com/gsabadini/go-bank-transfer/infrastructure"
	"github.com/gsabadini/go-bank-transfer/infrastructure/database"
	"github.com/gsabadini/go-bank-transfer/infrastructure/logger"
	"github.com/gsabadini/go-bank-transfer/infrastructure/validator"
	"github.com/gsabadini/go-bank-transfer/infrastructure/web"
)

func main() {
	infrastructure.NewConfig().
		AppName(os.Getenv("APP_NAME")).
		Logger(logger.InstanceLogrusLogger).
		DbSQL(database.InstancePostgres).
		DbNoSQL(database.InstanceMongoDB).
		Validator(validator.InstanceGoPlayground).
		WebServerPort(os.Getenv("APP_PORT")).
		WebServer(web.InstanceGorillaMux).
		StartApp()
}
