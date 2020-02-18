package main

import (
	"github.com/gsabadini/go-stone/api"
)

func main() {
	var server = api.NewHTTPServer()

	server.Listen()
}
