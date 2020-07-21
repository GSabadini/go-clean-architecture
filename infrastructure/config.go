package infrastructure

import (
	"os"
	"strconv"

	"github.com/gsabadini/go-bank-transfer/infrastructure/database"
	"github.com/gsabadini/go-bank-transfer/infrastructure/logger"
	"github.com/gsabadini/go-bank-transfer/infrastructure/validator"
	"github.com/gsabadini/go-bank-transfer/infrastructure/web"
	"github.com/gsabadini/go-bank-transfer/repository"
)

//Config armazena a estrutura de configuração da aplicação
type Config struct {
	appName   string
	port      web.Port
	WebServer web.Server
	Logger    logger.Logger
	dbSQL     repository.SQLHandler
	dbNoSQL   repository.NoSQLHandler
	validator validator.Validator
}

//NewConfig configura a aplicação
func NewConfig() Config {
	port, err := strconv.ParseInt(os.Getenv("APP_PORT"), 10, 64)
	if err != nil {
		panic(err)
	}

	config := Config{
		appName: os.Getenv("APP_NAME"),
		port:    web.Port(port),
		Logger:  log(),
	}

	config.validator = validation(config.Logger)
	config.dbSQL = postgresConn(config.Logger)
	config.dbNoSQL = mongoDBConn(config.Logger)

	config.WebServer = webServer(
		config.Logger,
		config.dbSQL,
		config.dbNoSQL,
		config.validator,
		config.port,
	)

	return config
}

func validation(log logger.Logger) validator.Validator {
	v, err := validator.NewValidatorFactory(validator.InstanceGoPlayground, log)
	if err != nil {
		panic(err)
	}

	log.Infof("Successfully configured validator")

	return v
}

func webServer(
	log logger.Logger,
	dbSQL repository.SQLHandler,
	dbNoSQL repository.NoSQLHandler,
	validator validator.Validator,
	port web.Port,
) web.Server {
	server, err := web.NewWebServerFactory(
		web.InstanceGorillaMux,
		log,
		dbSQL,
		dbNoSQL,
		validator,
		port,
	)

	if err != nil {
		panic(err)
	}

	log.Infof("Successfully configured web server")

	return server
}

func log() logger.Logger {
	log, err := logger.NewLoggerFactory(logger.InstanceLogrusLogger, true)
	if err != nil {
		panic(err)
	}

	log.Infof("Successfully configured logger")

	return log
}

func mongoDBConn(log logger.Logger) repository.NoSQLHandler {
	handler, err := database.NewDatabaseNoSQLFactory(
		database.InstanceMongoDB,
		database.NewConfigMongoDB(),
	)
	if err != nil {
		log.Fatalln("Could not make a connection to the database")
		panic(err)
	}

	log.Infof("Successfully connected to the NoSQL database")

	return handler
}

func postgresConn(log logger.Logger) repository.SQLHandler {
	handler, err := database.NewDatabaseSQLFactory(
		database.InstancePostgres,
		database.NewConfigPostgres(),
	)
	if err != nil {
		log.Fatalln("Could not make a connection to the database")
		panic(err)
	}

	log.Infof("Successfully connected to the SQL database")

	return handler
}
