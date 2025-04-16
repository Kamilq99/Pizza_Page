package main

import (
	"login_service/handlers"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	r.POST("/login", handlers.LoginUser)
	r.POST("/register", handlers.RegisterUser)

	r.Run(":8080")
}
