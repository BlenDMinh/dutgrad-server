package server

import (
	"strconv"

	"github.com/BlenDMinh/dutgrad-server/config"
)

func Init() {
	config := config.GetEnv()
	r := GetRouter()
	r.Run(":" + strconv.Itoa(config.Port))
}
