package controllers

import (
	"context"
	"log"

	"github.com/valerianomacuri/task-manager/common"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// Struct used for maintaining HTTP Request Context
type Context struct {
	mongo *mongo.Client
}

// Close mongo.Client
func (c *Context) Close() {
	err := c.mongo.Disconnect(context.TODO())
	if err != nil {
		log.Println("failed to disconnect")
		return
	}
	log.Println("disconnected from MongoDB")
}

// Returns mongo.collection for the given name
func (c *Context) DBCollection(name string) *mongo.Collection {
	return c.mongo.Database(common.AppConfig.Database).Collection(name)
}

// Create a new Context object for each HTTP request
func NewContext() *Context {
	ctx := context.TODO()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(common.AppConfig.MongoURI))
	if err = client.Ping(ctx, readpref.Primary()); err != nil {
		log.Fatal(err)
	}
	log.Println("Connected to MongoDB")
	return &Context{
		mongo: client,
	}
}
