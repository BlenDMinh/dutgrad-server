package main

import (
	"github.com/BlenDMinh/dutgrad-server/config"
	"github.com/BlenDMinh/dutgrad-server/server"
)

func main() {
	config.Init()
	server.Init()
}
