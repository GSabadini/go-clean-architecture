package infrastructure

import (
	"strconv"
	"time"

	"github.com/gsabadini/go-bank-transfer/infrastructure/database"
	"github.com/gsabadini/go-bank-transfer/infrastructure/logger"
	"github.com/gsabadini/go-bank-transfer/infrastructure/validator"
	"github.com/gsabadini/go-bank-transfer/infrastructure/web"
	"github.com/gsabadini/go-bank-transfer/repository"
)

//config armazena a estrutura de configuração da aplicação
type config struct {
	appName       string
	logger        logger.Logger
	validator     validator.Validator
	dbSQL         repository.SQLHandler
	dbNoSQL       repository.NoSQLHandler
	ctxTimeout    time.Duration
	webServerPort web.Port
	webServer     web.Server
}

//NewConfig configura a aplicação
func NewConfig() *config {
	return &config{}
}

func (c *config) Name(name string) *config {
	c.appName = name
	return c
}

func (c *config) Logger(instance int) *config {
	log, err := logger.NewLoggerFactory(instance)
	if err != nil {
		panic(err)
	}

	c.logger = log
	c.logger.Infof("Successfully configured log")
	return c
}

func (c *config) DbSQL(instance int) *config {
	db, err := database.NewDatabaseSQLFactory(instance)
	if err != nil {
		c.logger.Fatalln("Could not make a connection to the database")
		panic(err)
	}

	c.logger.Infof("Successfully connected to the SQL database")

	c.dbSQL = db
	return c
}

func (c *config) DbNoSQL(instance int) *config {
	db, err := database.NewDatabaseNoSQLFactory(instance)
	if err != nil {
		c.logger.Fatalln("Could not make a connection to the database")
		panic(err)
	}

	c.logger.Infof("Successfully connected to the NoSQL database")

	c.dbNoSQL = db
	return c
}

func (c *config) ContextTimeout(t time.Duration) *config {
	c.ctxTimeout = t
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
	s, err := web.NewWebServerFactory(
		instance,
		c.logger,
		c.dbSQL,
		c.dbNoSQL,
		c.validator,
		c.webServerPort,
		c.ctxTimeout,
	)

	if err != nil {
		panic(err)
	}

	c.logger.Infof("Successfully configured web server")

	c.webServer = s
	return c
}

func (c *config) WebServerPort(port string) *config {
	p, err := strconv.ParseInt(port, 10, 64)
	if err != nil {
		panic(err)
	}

	c.webServerPort = web.Port(p)
	return c
}

func (c *config) Start() {
	c.webServer.Listen()
}
