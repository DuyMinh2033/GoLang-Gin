package main

import (
	"ProjectGo/routes"
	"ProjectGo/utils"

	"github.com/gin-gonic/gin"
)

func main() {
	utils.InitMongoDB()
	r := gin.Default()
	r.Use(utils.SetupCORS())
	routes.RegisterRoutes(r)
	r.Run(":9000")
}
