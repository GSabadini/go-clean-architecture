package main

import (
	"github.com/gsabadini/go-bank-transfer/infrastructure"
)

func main() {
	var appConfig = infrastructure.NewConfig()

	appConfig.WebServer.Listen()
}
