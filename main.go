package main

import (
	"github.com/gsabadini/go-bank-transfer/api"
	"github.com/gsabadini/go-bank-transfer/config"
)

func main() {
	var (
		appConfig = config.NewConfig()
		server    = api.NewHTTPServer(appConfig)
	)

	server.Listen()
}
