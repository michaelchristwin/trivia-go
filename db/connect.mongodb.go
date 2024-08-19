package db

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var MongoClient *mongo.Client

func ConnectDB() {
	err := godotenv.Load(".env.local")
	if err != nil {
		log.Fatalf("Error loading .env.local file")
	}
	MONGO_PASSWORD := os.Getenv("ACCOUNT_PASSWORD")
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	uri := fmt.Sprintf("mongodb+srv://kelechukwuchristwin:%s@michael.fqimwas.mongodb.net/?retryWrites=true&w=majority&appName=michael", MONGO_PASSWORD)
	opts := options.Client().ApplyURI(uri).SetServerAPIOptions(serverAPI)
	MongoClient, err := mongo.Connect(context.TODO(), opts)
	if err != nil {
		panic(err)
	}
	defer func() {
		if err = MongoClient.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()
	if err := MongoClient.Database("admin").RunCommand(context.TODO(), bson.D{{Key: "Ping", Value: 1}}).Err(); err != nil {
		panic(err)
	}
	fmt.Println("Pinged your deployment. You successfully connected to MongoDB!")
}
