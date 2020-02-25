package config

import "os"

//Config armazena a estrutura de configuração da aplicação
type Config struct {
	AppName string
	ApiPort      int
	DatabaseName string
	DatabaseHost string
}

//NewConfig retorna a configuração da aplicação
func NewConfig() Config {
	return Config{
		AppName: "go-bank-transfer",
		ApiPort:      3001,
		DatabaseName: getDatabaseName(),
		DatabaseHost: getDatabaseHost(),
	}
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
