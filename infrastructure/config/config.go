package config

import (
	"fmt"
	"os"

	"github.com/gsabadini/go-bank-transfer/infrastructure/database"
	"github.com/gsabadini/go-bank-transfer/infrastructure/logger"
)

//Config armazena a estrutura de configuração da aplicação
type Config struct {
	AppName                 string
	APIPort                 int
	Logger                  logger.Logger
	DatabaseSQLConnection   database.SQLHandler
	DatabaseNOSQLConnection database.NoSQLHandler
}

//NewConfig retorna a configuração da aplicação
func NewConfig() Config {
	return Config{
		AppName:                 os.Getenv("APP_NAME"),
		APIPort:                 3001,
		Logger:                  getLogger(),
		DatabaseSQLConnection:   getConnectionDatabaseSQL(getLogger()),
		DatabaseNOSQLConnection: getConnectionDatabaseNoSQL(getLogger()),
	}
}

func getLogger() logger.Logger {
	log, err := logger.NewLogger(logger.Info, true, logger.InstanceZapLogger)
	if err != nil {
		panic(err)
	}

	log.Infof("Successfully configured logger")
	return log
}

func getConnectionDatabaseNoSQL(logger logger.Logger) database.NoSQLHandler {
	var (
		host   = verifyExistEnvironmentParams("MONGODB_HOST")
		dbName = verifyExistEnvironmentParams("MONGODB_DATABASE")
	)

	handler, err := database.NewDatabaseNoSQL(database.InstanceMongoDB, host, dbName)
	if err != nil {
		logger.Fatalln("Could not make a connection to the database")
		panic(err)
	}

	logger.Infof("Successfully connected to the database")

	return handler
}

func getConnectionDatabaseSQL(logger logger.Logger) database.SQLHandler {
	var ds = fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable password=%s",
		verifyExistEnvironmentParams("POSTGRES_HOST"),
		verifyExistEnvironmentParams("POSTGRES_PORT"),
		verifyExistEnvironmentParams("POSTGRES_USER"),
		verifyExistEnvironmentParams("POSTGRES_DATABASE"),
		verifyExistEnvironmentParams("POSTGRES_PASSWORD"),
	)
	handler, err := database.NewDatabaseSQL(database.InstancePostgres, ds)
	if err != nil {
		logger.Fatalln("Could not make a connection to the database")
		panic(err)
	}

	logger.Infof("Successfully connected to the database")

	return handler
}

func verifyExistEnvironmentParams(parameter string) string {
	if value := os.Getenv(parameter); value != "" {
		return value
	}

	panic(fmt.Sprintf("Environment variable '%s' has not been defined", parameter))
}
