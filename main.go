package main

import (
	"github.com/gsabadini/go-bank-transfer/api"
)

func main() {
	var server = api.NewHTTPServer()

	server.Listen()
}
