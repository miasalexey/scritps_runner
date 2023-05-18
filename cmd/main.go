package main

import (
	"github.com/gin-gonic/gin"
	"scritps_runner/internal/config"
	"scritps_runner/internal/utils"
)

func main() {
	r := gin.Default()
	utils.RegisterAllHandlers(r)
	cfg := config.GetConfig()
	err := r.Run(":" + cfg.ServerPort)
	if err != nil {
		return
	}
}
