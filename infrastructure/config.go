package infrastructure

import (
	"strconv"

	"github.com/gsabadini/go-bank-transfer/infrastructure/database"
	"github.com/gsabadini/go-bank-transfer/infrastructure/logger"
	"github.com/gsabadini/go-bank-transfer/infrastructure/validator"
	"github.com/gsabadini/go-bank-transfer/infrastructure/web"
	"github.com/gsabadini/go-bank-transfer/repository"
)

//config armazena a estrutura de configuração da aplicação
type config struct {
	appName       string
	webServerPort web.Port
	webServer     web.Server
	logger        logger.Logger
	dbSQL         repository.SQLHandler
	dbNoSQL       repository.NoSQLHandler
	validator     validator.Validator
}

//NewConfig configura a aplicação
func NewConfig() *config {
	return &config{}
}

func (c *config) AppName(n string) *config {
	c.appName = n
	return c
}

func (c *config) WebServerPort(p string) *config {
	port, err := strconv.ParseInt(p, 10, 64)
	if err != nil {
		panic(err)
	}

	c.webServerPort = web.Port(port)
	return c
}

func (c *config) Logger(instance int) *config {
	log, err := logger.NewLoggerFactory(instance, true)
	if err != nil {
		panic(err)
	}

	c.logger = log
	c.logger.Infof("Successfully configured log")
	return c
}

func (c *config) DbSQL(instance int) *config {
	handler, err := database.NewDatabaseSQLFactory(instance)
	if err != nil {
		c.logger.Fatalln("Could not make a connection to the database")
		panic(err)
	}

	c.logger.Infof("Successfully connected to the SQL database")

	c.dbSQL = handler
	return c
}

func (c *config) DbNoSQL(instance int) *config {
	handler, err := database.NewDatabaseNoSQLFactory(instance)
	if err != nil {
		c.logger.Fatalln("Could not make a connection to the database")
		panic(err)
	}

	c.logger.Infof("Successfully connected to the NoSQL database")

	c.dbNoSQL = handler
	return c
}

func (c *config) Validator(instance int) *config {
	v, err := validator.NewValidatorFactory(instance)
	if err != nil {
		panic(err)
	}

	c.logger.Infof("Successfully configured validator")

	c.validator = v
	return c
}

func (c *config) WebServer(instance int) *config {
	server, err := web.NewWebServerFactory(
		instance,
		c.logger,
		c.dbSQL,
		c.dbNoSQL,
		c.validator,
		c.webServerPort,
	)

	if err != nil {
		panic(err)
	}

	c.logger.Infof("Successfully configured web server")

	c.webServer = server
	return c
}

func (c *config) StartApp() {
	c.webServer.Listen()
}
