package utils

import (
	"github.com/gin-gonic/gin"
	"scritps_runner/internal/keepass_parser"
	"scritps_runner/internal/mikrotik"
)

func RegisterAllHandlers(router *gin.Engine) {
	mikrotikHandler := mikrotik.NewHandler()
	keePassParserHandler := keepass_parser.NewHandler()
	mikrotikHandler.Register(router)
	keePassParserHandler.Register(router)
}
