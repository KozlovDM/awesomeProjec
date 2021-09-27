package main

import (
	controller "awesomeProject/client/restController"
	"github.com/gin-gonic/gin"
	"log"
)

func main() {
	restServer := gin.Default()

	restServer.POST("/signup", controller.SignUp)
	restServer.POST("/signin", controller.SignIn)

	if err := restServer.Run(); err != nil {
		log.Fatalf("failed to serve: %s", err)
	}
}
