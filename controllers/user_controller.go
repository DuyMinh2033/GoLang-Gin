package controllers

import (
	models "ProjectGo/model"
	"ProjectGo/utils"
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	UserCollection *mongo.Collection = utils.GetCollection("users")
	jobQueue                         = make(chan models.User, 100)
)

func init() {
	fmt.Println("Initializing worker pool...")
	for i := 0; i < 5; i++ {
		go worker()
	}
}

func worker() {
	for user := range jobQueue {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		_, err := UserCollection.InsertOne(ctx, user)
		if err != nil {
			fmt.Printf("Failed to insert user: %v\n", err)
		} else {
			fmt.Println("User inserted successfully:", user)
		}
		cancel()
	}
}

func CreateUserHandler(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON payload"})
		return
	}
	jobQueue <- user
	c.JSON(http.StatusOK, gin.H{
		"message": "User creation request received",
		"user":    user,
	})
}
