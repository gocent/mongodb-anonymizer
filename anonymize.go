package main

import (
	"context"
	"encoding/json"
	"fmt"
	"gocent/mongodb-anonymizer/config"

	"github.com/go-faker/faker/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
) // import the faker package

func main() {
	fmt.Println("Hello " + faker.Name())
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(config.GetEnv().DB.SourceURI.String()))
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
	fmt.Println("Collections: ", config.GetEnv().DB.Collections)
	for _, collection := range config.GetEnv().DB.Collections {
		cursor, err := client.Database(config.GetEnv().DB.SourceName).Collection(collection).Find(context.TODO(), bson.M{})
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
