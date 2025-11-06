package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Server struct {
}

func (Server) GetHello(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Hello from oapi-codegen with Gin!",
	})
}
