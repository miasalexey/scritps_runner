package main

import (
	"github.com/gin-gonic/gin"
	"scritps_runner/internal/mikrotik"
)

func main() {
	r := gin.Default()
	m := mikrotik.NewHandler()
	m.Register(r)
	err := r.Run(":4003")
	if err != nil {
		return
	}
}
