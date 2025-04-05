package controllers

import (
	models "ProjectGo/model"
	"ProjectGo/utils"
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	TodoCollection *mongo.Collection = utils.GetCollection("todos")
	todoJobQueue                     = make(chan models.Todo, 100)
	deleteJobQueue                   = make(chan primitive.ObjectID, 100)
)

func init() {
	fmt.Println("Initializing todo worker pool...")
	for i := 0; i < 5; i++ {
		go insertTodo()
		go deleteTodo()
	}
}

func insertTodo() {
	for todo := range todoJobQueue {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		_, err := TodoCollection.InsertOne(ctx, todo)
		if err != nil {
			fmt.Printf("Failed to insert todo: %v\n", err)
		} else {
			fmt.Println("Todo inserted successfully:", todo)
		}
		cancel()
	}
}

func deleteTodo() {
	for todoID := range deleteJobQueue {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second) // Xóa todo từ MongoDB
		filter := bson.M{"_id": todoID}
		result, err := TodoCollection.DeleteOne(ctx, filter)
		if err != nil {
			fmt.Printf("Failed to delete todo: %v\n", err)
		} else {
			fmt.Printf("Todo deleted successfully with ID: %v\n", todoID)
		}
		if result.DeletedCount == 0 {
			fmt.Println("Todo not found to delete:", todoID)
		}
		cancel()
	}
}

func CreateTodoHandler(c *gin.Context) {
	var todo models.Todo
	if err := c.ShouldBindJSON(&todo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON payload"})
		return
	}
	todo.ID = primitive.NewObjectID()
	todoJobQueue <- todo
	c.JSON(http.StatusOK, gin.H{
		"message": "Todo creation request received",
		"todo":    todo,
	})
}

func GetAllTodosHandler(c *gin.Context) {
	var todos []models.Todo

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := TodoCollection.Find(ctx, bson.M{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch todos"})
		return
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var todo models.Todo
		if err := cursor.Decode(&todo); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode todo"})
			return
		}
		todos = append(todos, todo)
	}
	if len(todos) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "No todos found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Todos fetched successfully",
		"todos":   todos,
	})
}

func DeleteTodoHandler(c *gin.Context) {
	id := c.Param("id")
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}
	deleteJobQueue <- objectID
	c.JSON(http.StatusOK, gin.H{
		"message": "Todo deletion request received",
		"id":      id,
	})
}
