package main

import (
	"github.com/BlenDMinh/dutgrad-server/configs"
	"github.com/BlenDMinh/dutgrad-server/databases"
	"github.com/BlenDMinh/dutgrad-server/server"
)

func main() {
	configs.Init()
	databases.Init()
	server.Init()
}
