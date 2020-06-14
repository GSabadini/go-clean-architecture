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
		AppName:                 "go-bank-transfer",
		APIPort:                 3001,
		Logger:                  getLogger(),
		DatabaseSQLConnection:   getConnectionPostgres(getLogger()),
		DatabaseNOSQLConnection: getConnectionMongoDB(getLogger()),
	}
}

func getLogger() logger.Logger {
	log, err := logger.NewLogger(logger.InstanceLogrusLogger)
	if err != nil {
		panic(err)
	}

	log.Infof("Successfully configured logger")
	return log
}

func getConnectionMongoDB(logger logger.Logger) *database.MongoHandler {
	handler, err := database.NewMongoHandler(
		verifyExistEnvironmentParams("MONGODB_HOST"),
		verifyExistEnvironmentParams("MONGODB_DATABASE"),
	)

	if err != nil {
		logger.Fatalln("Could not make a connection to the database")
		panic(err)
	}

	logger.Infof("Successfully connected to the database")

	return handler
}

func getConnectionPostgres(logger logger.Logger) *database.PostgresHandler {
	ds := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable password=%s",
		verifyExistEnvironmentParams("POSTGRES_HOST"),
		verifyExistEnvironmentParams("POSTGRES_PORT"),
		verifyExistEnvironmentParams("POSTGRES_USER"),
		verifyExistEnvironmentParams("POSTGRES_DATABASE"),
		verifyExistEnvironmentParams("POSTGRES_PASSWORD"),
	)

	handler, err := database.NewPostgresHandler(ds)
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
