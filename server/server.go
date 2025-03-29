package server

import (
	"strconv"

	"github.com/BlenDMinh/dutgrad-server/configs"
)

func Init() {
	config := configs.GetEnv()
	r := GetRouter()
	r.Run(":" + strconv.Itoa(config.Port))
}

func Close() {

}
