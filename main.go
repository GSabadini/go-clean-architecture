package main

import (
	"github.com/gsabadini/go-bank-transfer/infrastructure"
)

func main() {
	var config = infrastructure.NewConfig()
	config.WebServer.Listen()
}
