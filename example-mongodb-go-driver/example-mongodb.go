package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/mongodb/mongo-go-driver/bson/objectid"
	"github.com/mongodb/mongo-go-driver/mongo"
)

// This is a type to hold our word definitions in
// we specifiy both bson (for MongoDB) and json (for web)
// naming for marshalling and unmarshalling
type item struct {
	ID         objectid.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	Word       string            `bson:"word" json:"word"`
	Definition string            `bson:"definition" json:"definition"`
}

var client *mongo.Client

func wordHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		// Use the session to get the DB and then the collection
		// Using an empty string for DB() gets the datbase specified
		// in the connection string
		c := client.Database("grand_tour").Collection("words")

		// Create an array of
		var items []item
		cur, err := c.Find(context.Background(), nil)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer cur.Close(context.Background())

		for cur.Next(context.Background()) {
			item := item{}
			err := cur.Decode(&item)
			if err != nil {
				log.Fatal("Decode check ", err)
			}
			items = append(items, item)
		}
		if err := cur.Err(); err != nil {
			log.Fatal("Cursor check ", err)
		}

		jsonstr, err := json.Marshal(items)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonstr)
		return
	case "PUT":
		r.ParseForm()
		c := client.Database("grand_tour").Collection("words")
		newItem := item{ID: objectid.New(), Word: r.Form.Get("word"), Definition: r.Form.Get("definition")}
		_, err := c.InsertOne(context.Background(), newItem, nil)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusAccepted)
		return
	}

	return
}

func main() {
	// Connect to database:
	// Connection string in $COMPOSE_MONGODB_URL
	// Compose database certificate pointed to in $PATH_TO_MONGODB_CERT

	// Get the environment variables with the connection string and certpath

	connectionString, present := os.LookupEnv("COMPOSE_MONGODB_URL")
	if !present {
		log.Fatal("Need to set COMPOSE_MONGODB_URL environment variable")
	}

	certpath, certavail := os.LookupEnv("PATH_TO_MONGODB_CERT")

	var err error

	if certavail {
		client, err = mongo.NewClientWithOptions(connectionString, mongo.ClientOpt.SSLCaFile(certpath))
	} else {
		client, err = mongo.NewClient(connectionString)
	}
	if err != nil {
		log.Fatal(err)
	}

	err = client.Connect(context.TODO())

	if err != nil {
		log.Fatal(err)
	}

	defer client.Disconnect(context.TODO())

	fs := http.FileServer(http.Dir("public"))
	http.Handle("/", fs)
	http.HandleFunc("/words", wordHandler)
	fmt.Println("Listening on localhost:8080")
	http.ListenAndServe(":8080", nil)
}
