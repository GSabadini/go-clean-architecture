package config

import (
	"os"

	"github.com/sirupsen/logrus"
)

//Config armazena a estrutura de configuração da aplicação
type Config struct {
	AppName      string
	ApiPort      int
	Logger       *logrus.Logger
	DatabaseName string
	DatabaseHost string
}

//NewConfig retorna a configuração da aplicação
func NewConfig() Config {
	return Config{
		AppName:      "go-bank-transfer",
		ApiPort:      3001,
		Logger:       createLoggerApp(),
		DatabaseName: getDatabaseName(),
		DatabaseHost: getDatabaseHost(),
	}
}

//TODO REVER LOCAL DE CRIAÇÃO DO LOGGER
func createLoggerApp() *logrus.Logger {
	log := logrus.New()
	log.SetFormatter(&logrus.JSONFormatter{})
	return log
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
