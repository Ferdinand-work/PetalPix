package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/Ferdinand-work/PetalPix/controllers"
	"github.com/Ferdinand-work/PetalPix/services"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var (
	server         *gin.Engine
	userService    services.UserService
	UserController *controllers.UserController
	ctx            context.Context
	userCollection *mongo.Collection
	mongoclient    *mongo.Client
	err            error
)

func init() {
	err = godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
		return
	}

	mongoURI := os.Getenv("MONGO_URI")
	if mongoURI == "" {
		log.Fatal("MONGO_URI environment variable not set")
	}

	ctx = context.TODO()
	mongoconn := options.Client().ApplyURI(mongoURI)
	mongoclient, err = mongo.Connect(ctx, mongoconn)
	if err != nil {
		log.Fatal(err)
	}
	err = mongoclient.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Mongo Connection Established")
	userCollection = mongoclient.Database("petalPix").Collection("users")
	userService = services.NewUserService(userCollection, ctx)
	UserController = controllers.New(userService)
	server = gin.Default()
}

func main() {
	defer mongoclient.Disconnect(ctx)
	basepath := server.Group("/v1")
	UserController.RegisterUserRoutes(basepath)
	Port := os.Getenv("PORT")
	if Port == "" {
		log.Fatal("PORT environment variable not set")
	}
	log.Fatal(server.Run(":" + Port))
}
