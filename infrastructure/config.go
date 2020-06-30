package infrastructure

import (
	"fmt"
	"os"
	"strconv"

	"github.com/gsabadini/go-bank-transfer/infrastructure/database"
	"github.com/gsabadini/go-bank-transfer/infrastructure/logger"
	"github.com/gsabadini/go-bank-transfer/infrastructure/validator"
	"github.com/gsabadini/go-bank-transfer/infrastructure/web"
)

//Config armazena a estrutura de configuração da aplicação
type Config struct {
	appName   string
	port      web.Port
	WebServer web.Server
	Logger    logger.Logger
	dbSQL     database.SQLHandler
	dbNoSQL   database.NoSQLHandler
	validator validator.Validator
}

//NewConfig retorna a configuração da aplicação
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
	config.dbSQL = dbSQLConn(config.Logger)
	config.dbNoSQL = dbNoSQLConn(config.Logger)

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
	v, err := validator.NewValidator(validator.InstanceGoPlayground, log)
	if err != nil {
		panic(err)
	}

	log.Infof("Successfully configured validator")

	return v
}

func webServer(
	log logger.Logger,
	dbSQLConn database.SQLHandler,
	dbNoSQLConn database.NoSQLHandler,
	validator validator.Validator,
	port web.Port,
) web.Server {
	server, err := web.NewWebServer(
		web.InstanceGorillaMux,
		log,
		dbSQLConn,
		dbNoSQLConn,
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
	log, err := logger.NewLogger(logger.InstanceLogrusLogger, true)
	if err != nil {
		panic(err)
	}

	log.Infof("Successfully configured logger")

	return log
}

func dbNoSQLConn(log logger.Logger) database.NoSQLHandler {
	var (
		host   = verifyExistEnvironmentParams("MONGODB_HOST")
		dbName = verifyExistEnvironmentParams("MONGODB_DATABASE")
	)

	handler, err := database.NewDatabaseNoSQL(database.InstanceMongoDB, host, dbName)
	if err != nil {
		log.Fatalln("Could not make a connection to the database")
		panic(err)
	}

	log.Infof("Successfully connected to the database")

	return handler
}

func dbSQLConn(log logger.Logger) database.SQLHandler {
	var ds = fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable password=%s",
		verifyExistEnvironmentParams("POSTGRES_HOST"),
		verifyExistEnvironmentParams("POSTGRES_PORT"),
		verifyExistEnvironmentParams("POSTGRES_USER"),
		verifyExistEnvironmentParams("POSTGRES_DATABASE"),
		verifyExistEnvironmentParams("POSTGRES_PASSWORD"),
	)

	handler, err := database.NewDatabaseSQL(database.InstancePostgres, ds)
	if err != nil {
		log.Fatalln("Could not make a connection to the database")
		panic(err)
	}

	log.Infof("Successfully connected to the database")

	return handler
}

func verifyExistEnvironmentParams(parameter string) string {
	if value := os.Getenv(parameter); value != "" {
		return value
	}

	panic(fmt.Sprintf("Environment variable '%s' has not been defined", parameter))
}
