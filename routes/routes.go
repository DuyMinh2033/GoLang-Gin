package routes

import (
	"ProjectGo/controllers"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine) {
	// Route cho User
	r.GET("/healthy", controllers.TestApiWorking)
	r.POST("/create", controllers.CreateUserHandler)
}
