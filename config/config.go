package config

import (
	"os"

	"github.com/gsabadini/go-bank-transfer/infrastructure/database"
	"github.com/sirupsen/logrus"
)

//Config armazena a estrutura de configuração da aplicação
type Config struct {
	AppName            string
	APIPort            int
	Logger             *logrus.Logger
	DatabaseConnection *database.MongoHandler
	DatabaseName       string
	DatabaseHost       string
}

//NewConfig retorna a configuração da aplicação
func NewConfig() Config {
	return Config{
		AppName:            "go-bank-transfer",
		APIPort:            3001,
		Logger:             getLogger(),
		DatabaseName:       getDatabaseName(),
		DatabaseHost:       getDatabaseHost(),
		DatabaseConnection: getDatabaseConnection(getLogger()),
	}
}

func getLogger() *logrus.Logger {
	log := logrus.New()
	log.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
	})

	return log
}

func getDatabaseConnection(logger *logrus.Logger) *database.MongoHandler {
	handler, err := database.NewMongoHandler(getDatabaseHost(), getDatabaseName())
	if err != nil {
		logger.Infoln("Could not make a connection to the database")

		// Se não conseguir conexão com o banco por algum motivo, então a aplicação deve criticar
		panic(err)
	}

	logger.Infoln("Successfully connected to the database")

	return handler
}

func getDatabaseHost() string {
	if host := os.Getenv("MONGODB_HOST"); host != "" {
		return host
	}

	panic("Environment variable 'MONGODB_HOST' has not been defined")
}

func getDatabaseName() string {
	if name := os.Getenv("MONGODB_DATABASE"); name != "" {
		return name
	}

	panic("Environment variable 'MONGODB_DATABASE' has not been defined")
}
