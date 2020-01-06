package main

import (
	"os"
	"context"
	"log"
	"strconv"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Trainer struct {
	Name string
	Age  int
	City string
}

func main() {
	log.Println("Initiation of Mongo Connection")
	MONGO_HOST := os.Getenv("MONGO_HOST")
	MONGO_PORT := os.Getenv("MONGO_PORT")

	clientOptions := options.Client().ApplyURI("mongodb://" + MONGO_HOST + ":" + MONGO_PORT)

	client, err := mongo.Connect(context.TODO(), clientOptions)

	if err != nil{
		log.Fatal(err)
	}

	err = client.Ping(context.TODO(), nil)

	if err != nil{
		log.Fatal(err)
	}

	log.Println("Connected to MongoDb")

	collection := client.Database("test").Collection("trainers")

	// Inserting Single Document
	// --------------------------
	ash := Trainer{"Ash", 10, "Pallet Town"}
	misty := Trainer{"Misty", 10, "Cerulean City"}
	brock := Trainer{"Brock", 15, "Power City"}

	insertResult, err := collection.InsertOne(context.TODO(), ash)
	if err != nil{
		log.Fatal(err)
	}
	log.Println("Inserted a single document: ", insertResult.InsertedID)

	// Inserting Multiple Documents
	// ----------------------------------
	trainers := []interface{}{misty, brock}
	insertManyResult, err := collection.InsertMany(context.TODO(), trainers)

	if err != nil{
		log.Fatal(err)
	}
	log.Println("Inserted a Multiple document: ", insertManyResult.InsertedIDs)

	// Updating Documents
	// --------------------------------
	filter := bson.D{{"name", "Ash"}}

	update := bson.D{
		{"$inc", bson.D{
			{"age", 1},
		}},
	}

	updateResult, err := collection.UpdateOne(context.TODO(), filter, update)
	if err != nil{
		log.Fatal(err)
	}

	// convert int64 to string
	matchCount := strconv.FormatInt(updateResult.MatchedCount, 10)
	modifiedCount := strconv.FormatInt(updateResult.ModifiedCount, 10)

	log.Println("Matched " + matchCount + " documents and updated " + modifiedCount + " documents. \n")

	// Finding a Single Document
	// ---------------------------
	var result *Trainer
	
	err = collection.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil{
		log.Fatal(err)
	}
	log.Println("Found a single document: " + result.Name + "\n")

	// Finding a Multple Documents 
	// ---------------------------
	findOptions := options.Find()
	findOptions.SetLimit(2)

	var results []*Trainer

	cur, err := collection.Find(context.TODO(), bson.D{{}}, findOptions)
	if err != nil{
		log.Fatal(err)
	}

	for cur.Next(context.TODO()) {
		var elem Trainer
		err := cur.Decode(&elem)
		if err != nil{
			log.Fatal(err)
		}

		results = append(results, &elem)
	}
	if err := cur.Err(); err != nil{
		log.Fatal(err)
	}

	cur.Close(context.TODO())
	foundCount := strconv.Itoa(len(results))
	log.Println("Found " + foundCount + " documents (array of pointers).")


	// Deleting documents
	// --------------------
	// deleteResult, err = collection.DeleteOne(context.TODO(), filter)
	deleteResult, err := collection.DeleteMany(context.TODO(), bson.D{{}})
	if err := cur.Err(); err != nil{
		log.Fatal(err)
	}
	deletedCount := strconv.FormatInt(deleteResult.DeletedCount, 10)
	log.Println("Deleted " + deletedCount + ".")

	// Closing the Mongodb connection
	// ----------------------------------
	err = client.Disconnect(context.TODO())

	if err != nil{
		log.Fatal(err)
	}
	log.Println("Connection to MongoDb closed.")
}