package config

import (
	"fmt"
	"os"

	"github.com/gsabadini/go-bank-transfer/infrastructure/database"
	"github.com/sirupsen/logrus"
)

//Config armazena a estrutura de configuração da aplicação
type Config struct {
	AppName            string
	APIPort            int
	Logger             *logrus.Logger
	DatabaseConnection *database.PostgresHandler
}

//NewConfig retorna a configuração da aplicação
func NewConfig() Config {
	return Config{
		AppName:            "go-bank-transfer",
		APIPort:            3001,
		Logger:             getLogger(),
		DatabaseConnection: getConnectionPostgres(getLogger()),
	}
}

func getLogger() *logrus.Logger {
	log := logrus.New()
	log.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
	})

	return log
}

func getConnectionMongoDB(logger *logrus.Logger) *database.MongoHandler {
	handler, err := database.NewMongoHandler(
		verifyExistEnvironmentParams("MONGODB_HOST"),
		verifyExistEnvironmentParams("MONGODB_DATABASE"),
	)

	if err != nil {
		logger.Infoln("Could not make a connection to the database")

		// Se não conseguir conexão com o banco por algum motivo, então a aplicação deve criticar
		panic(err)
	}

	logger.Infoln("Successfully connected to the database")

	return handler
}

func getConnectionPostgres(logger *logrus.Logger) *database.PostgresHandler {
	ds := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable password=%s",
		verifyExistEnvironmentParams("POSTGRES_HOST"),
		verifyExistEnvironmentParams("POSTGRES_PORT"),
		verifyExistEnvironmentParams("POSTGRES_USER"),
		verifyExistEnvironmentParams("POSTGRES_DATABASE"),
		verifyExistEnvironmentParams("POSTGRES_PASSWORD"),
	)

	handler, err := database.NewPostgresHandler(ds)
	if err != nil {
		logger.Infoln("Could not make a connection to the database")

		// Se não conseguir conexão com o banco por algum motivo, então a aplicação deve criticar
		panic(err)
	}

	logger.Infoln("Successfully connected to the database")

	return handler
}

func verifyExistEnvironmentParams(parameter string) string {
	if value := os.Getenv(parameter); value != "" {
		return value
	}

	panic(fmt.Sprintf("Environment variable '%s' has not been defined", parameter))
}
