package main

import (
	"ProjectGo/controllers"
	"ProjectGo/routes"
	"ProjectGo/utils"
	"context"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("don't load .env file")
	}
	uri := os.Getenv("MONGO_URI")
	client := utils.ConnectMongoDB(uri)
	defer func() {
		if err := client.Disconnect(context.Background()); err != nil {
			log.Fatalf("Don't close connect with MongoDB: %v", err)
		}
	}()
	controllers.UserCollection = client.Database("testdb").Collection("users")
	r := gin.Default()
	r.Use(utils.SetupCORS())
	routes.RegisterRoutes(r)
	r.Run(":9000")
}
