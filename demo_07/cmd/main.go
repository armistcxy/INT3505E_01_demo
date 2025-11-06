package main

import (
	"demo_07/api"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	server := api.Server{}
	api.RegisterHandlers(r, server)

	log.Println("Server started at :8080")
	r.Run(":8080")
}
