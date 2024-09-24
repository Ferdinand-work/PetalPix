package main

import (
	"context"
	"os"
	"time"

	"github.com/Ferdinand-work/PetalPix/controllers"
	"github.com/Ferdinand-work/PetalPix/services"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"gopkg.in/yaml.v3"
)

var (
	server         *gin.Engine
	userService    services.UserService
	UserController *controllers.UserController
	ctx            context.Context
	userCollection *mongo.Collection
	mongoclient    *mongo.Client
	Port           string
)

type Config struct {
	Environments map[string]Environment `yaml:"environments"`
}

// Environment holds the environment variables
type Environment struct {
	PORT       string `yaml:"PORT"`
	MONGO_URI  string `yaml:"MONGO_URI"`
	DB         string `yaml:"DB"`
	COLLECTION string `yaml:"COLLECTION"`
}

func logMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		logrus.Infof("Request: %s %s | Status: %d | Duration: %v", c.Request.Method, c.Request.URL.Path, c.Writer.Status(), time.Since(start))
	}
}

func init() {
	// Set up logging to file
	logFile, err := os.OpenFile("./logs/app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		logrus.Fatalf("Failed to open log file: %v", err)
	}
	logrus.SetOutput(logFile)
	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})

	// Read configuration file
	data, err := os.ReadFile("config.yml")
	if err != nil {
		logrus.Fatalf("Failed to read config file: %v", err)
	}

	var config Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		logrus.Fatalf("Failed to unmarshal config: %v", err)
	}

	// Get environment settings
	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "local" // Default to local if no environment variable is set
	}
	envConfig := config.Environments[env]
	mongoURI := envConfig.MONGO_URI
	DB := envConfig.DB
	COLLECTION := envConfig.COLLECTION
	Port = envConfig.PORT

	// Validate Mongo URI
	if mongoURI == "" {
		logrus.Fatal("MONGO_URI environment variable not set")
	}

	// Connect to MongoDB
	ctx = context.TODO()
	mongoconn := options.Client().ApplyURI(mongoURI)
	mongoclient, err = mongo.Connect(ctx, mongoconn)
	if err != nil {
		logrus.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	err = mongoclient.Ping(ctx, readpref.Primary())
	if err != nil {
		logrus.Fatalf("Failed to ping MongoDB: %v", err)
	}
	logrus.Info("Mongo connection established")
	userCollection = mongoclient.Database(DB).Collection(COLLECTION)
	userService = services.NewUserService(userCollection, ctx)
	UserController = controllers.New(userService)

	// Initialize Gin server
	server = gin.Default()

	server.Use(logMiddleware())

}

func main() {
	defer func() {
		if err := mongoclient.Disconnect(ctx); err != nil {
			logrus.Errorf("Failed to disconnect MongoDB client: %v", err)
		} else {
			logrus.Info("MongoDB client disconnected")
		}
	}()

	// Set up routes
	basepath := server.Group("/v1")
	UserController.RegisterUserRoutes(basepath)

	// Validate port and start server
	if Port == "" {
		logrus.Fatal("PORT environment variable not set")
	}
	logrus.Infof("Starting server on port %s", Port)
	if err := server.Run(":" + Port); err != nil {
		logrus.Fatalf("Failed to run server: %v", err)
	}
}
