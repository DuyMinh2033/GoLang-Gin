package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	ID   string `bson:"_id,omitempty"`
	Name string `json:"name" bson:"name"`
	Age  string `json:"age" bson:"age"`
}

type Todo struct {
	ID   primitive.ObjectID `json:"id" bson:"_id"`
	Todo string             `json:"todo" bson:"todo"`
}
