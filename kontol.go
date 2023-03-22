package main

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"os"
)

func main() {
	// send it to MONGODB
	// retrive it's configuration :<
	ctx := context.TODO()
	mongouri := os.Getenv("MONGODB_CONNECT_URI")
	database := "bulletin"
	collection := "posts"

	// and the job
	// Connect
	mclient, err := mongo.Connect(ctx, options.Client().ApplyURI(mongouri))

	if err != nil {
		log.Fatal(err)
	}

	type doctemplate struct {
		title        string `bson:"title,omitempty"`
		descriptive  string `bson:"descriptive,omitempty"`
		thumbnailuri string `bson:"thumbnailuri,omitempty"`
	}

	_, err = mclient.Database(database).Collection(collection).InsertOne(ctx, bson.M{
		"kontol": "china",
		"asu":    9121.12,
		"memek":  true,
		"jepang": nil,
	})

	if err != nil {
		log.Fatal(err)
	}

	defer func(mclient *mongo.Client, ctx context.Context) {
		err := mclient.Disconnect(ctx)
		if err != nil {
			panic(err)
		}
	}(mclient, ctx)
}
