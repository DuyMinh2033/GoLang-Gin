package routes

import (
	"ProjectGo/controllers"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine) {
	// Route cho User
	r.POST("/create", controllers.CreateUserHandler)
	r.POST("/addTodo", controllers.CreateTodoHandler)
	r.GET("/getAll", controllers.GetAllTodosHandler)
	r.DELETE("/deleteTodo/:id", controllers.DeleteTodoHandler)
}
