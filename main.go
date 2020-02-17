package main

import (
	"github.com/gsabadini/go-stone/api"
)

func main() {
	server := api.HTTPServer{}
	server.Listen()
}
