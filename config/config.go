package config

import "os"

//Config armazena a estrutura de configuração da aplicação
type Config struct {
	ApiPort      int
	DatabaseName string
	DatabaseHost string
}

//NewConfig retorna a configuração da aplicação
func NewConfig() Config {
	return Config{
		ApiPort:      3001,
		DatabaseName: getDatabaseName(),
		DatabaseHost: getDatabaseHost(),
	}
}

func getDatabaseHost() string {
	if uri := os.Getenv("MONGODB_HOST"); uri != "" {
		return uri
	}

	panic("Environment variable 'MONGODB_HOST' has not been defined")
}

func getDatabaseName() string {
	if uri := os.Getenv("MONGODB_DATABASE"); uri != "" {
		return uri
	}

	panic("Environment variable 'MONGODB_DATABASE' has not been defined")
}
