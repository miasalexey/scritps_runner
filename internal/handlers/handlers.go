package handlers

import (
	"github.com/gin-gonic/gin"
)

// Handlers interface for subsequent implementations by structures
type Handlers interface {
	Register(router *gin.Engine)
}
