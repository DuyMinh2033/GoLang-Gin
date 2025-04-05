package utils

import (
	"context"
	"log"
	"os"
	"sync"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func SetupCORS() gin.HandlerFunc {
	return cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:8080"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	})
}

var (
	clientInstance *mongo.Client
	clientOnce     sync.Once
	databaseName   = "testdb"
)

func InitMongoDB() *mongo.Client {
	clientOnce.Do(func() {
		err := godotenv.Load()
		if err != nil {
			log.Println("Warning: .env file not loaded")
		}

		uri := os.Getenv("MONGO_URI")
		if uri == "" {
			log.Fatal("MONGO_URI is not set in the environment variables")
		}

		clientOptions := options.Client().ApplyURI(uri)
		client, err := mongo.Connect(context.Background(), clientOptions)
		if err != nil {
			log.Fatalf("Failed to connect to MongoDB: %v", err)
		}

		// Ping MongoDB to verify connection
		err = client.Ping(context.Background(), nil)
		if err != nil {
			log.Fatalf("Failed to ping MongoDB: %v", err)
		}

		log.Println("Connected to MongoDB!")
		clientInstance = client
	})
	return clientInstance
}

// GetCollection returns a MongoDB collection by name
func GetCollection(collectionName string) *mongo.Collection {
	client := InitMongoDB()
	return client.Database(databaseName).Collection(collectionName)
}
