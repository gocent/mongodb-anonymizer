package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/go-faker/faker/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
) // import the faker package

const uri = "mongodb://127.0.0.1:27017"

func main() {
	fmt.Println("Hello " + faker.Name())
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		panic(err)
	}
	defer func() {
		if err = client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()
	// Ping the primary
	if err := client.Ping(context.TODO(), readpref.Primary()); err != nil {
		panic(err)
	}
	fmt.Println("Successfully connected and pinged.")

	cursor, err := client.Database("anonymize").Collection("test").Find(context.TODO(), bson.M{})

	if err != nil {
		panic(err)
	}

	var results []interface{}

	if err = cursor.All(context.TODO(), &results); err != nil {
		panic(err)
	}

	for _, result := range results {
		res, _ := json.Marshal(result)
		replaceAllProperties(res)
		fmt.Println(string(res))
	}
}

func replaceAllProperties(data interface{}) {
	switch data.(type) {
	case map[string]interface{}:
		for value := range data.(map[string]interface{}) {
			replaceAllProperties(value)
		}
	case []interface{}:
		for _, value := range data.([]interface{}) {
			replaceAllProperties(value)
		}
	default:
		fmt.Println(data)
	}
}
