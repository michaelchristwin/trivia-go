package db

import (
	"context"
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var MongoClient *mongo.Client

func ConnectDB() error {
	err := godotenv.Load(".env.local")
	if err != nil {
		return fmt.Errorf("error loading .env.local file: %w", err)
	}

	MONGO_PASSWORD := os.Getenv("ACCOUNT_PASSWORD")
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	uri := fmt.Sprintf("mongodb+srv://kelechukwuchristwin:%s@michael.fqimwas.mongodb.net/?retryWrites=true&w=majority&appName=michael", MONGO_PASSWORD)
	opts := options.Client().ApplyURI(uri).SetServerAPIOptions(serverAPI)
	client, err := mongo.Connect(context.TODO(), opts)
	if err != nil {
		return fmt.Errorf("error connecting to MongoDB: %w", err)
	}

	// Verify connection
	if err := client.Database("admin").RunCommand(context.TODO(), bson.D{{Key: "Ping", Value: 1}}).Err(); err != nil {
		return fmt.Errorf("error verifying MongoDB connection: %w", err)
	}

	MongoClient = client
	fmt.Println("Pinged your deployment. You successfully connected to MongoDB!")
	return nil
}
