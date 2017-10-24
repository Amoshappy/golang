package main

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"strings"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// This is a type to hold our word definitions in
// we specifiy both bson (for MongoDB) and json (for web)
// naming for marshalling and unmarshalling
type item struct {
	ID         bson.ObjectId `bson:"_id,omitempty" json:"_id,omitempty"`
	Word       string        `bson:"word" json:"word"`
	Definition string        `bson:"definition" json:"definition"`
}

var session *mgo.Session

func wordHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		// Use the session to get the DB and then the collection
		// Using an empty string for DB() gets the datbase specified
		// in the connection string
		c := session.DB("grand_tour").C("words")

		// Create an array of
		var items []item
		err := c.Find(nil).All(&items)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
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
		c := session.DB("grand_tour").C("words")
		newItem := item{Word: r.Form.Get("word"), Definition: r.Form.Get("definition")}
		err := c.Insert(newItem)
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

	// Building a TLS configuration
	// Create a certificate pool for root certificates
	// Then load that pool with our certificate
	roots := x509.NewCertPool()
	if ca, err := ioutil.ReadFile(os.Getenv("PATH_TO_MONGODB_CERT")); err == nil {
		roots.AppendCertsFromPEM(ca)
	}
	// Then create a TLS config object
	// and save that new root pool in it as it's only CA
	tlsConfig := &tls.Config{}
	tlsConfig.RootCAs = roots

	// Get the environment variable with the connection string
	connectionString := os.Getenv("COMPOSE_MONGODB_URL")
	// Currently mgo errors out if it sees the ?ssl=true option on the
	// connectionString, so we'll trim that off
	trimmedConnectionString := strings.TrimSuffix(connectionString, "?ssl=true")
	// Now we can parse the connection string
	dialInfo, err := mgo.ParseURL(trimmedConnectionString)
	if err != nil {
		log.Fatal(err)
	}

	// Ok, now we can modify the Dial process for the connection making
	// the DialServer function actually use a TLS version of Dial and
	// pass it our TLS configuration
	dialInfo.DialServer = func(addr *mgo.ServerAddr) (net.Conn, error) {
		return tls.Dial("tcp", addr.String(), tlsConfig)
	}

	// Using that modified DialInfo, we can now connect to the database
	session, err = mgo.DialWithInfo(dialInfo)
	if err != nil {
		log.Fatal(err)
	}

	defer session.Close()

	fs := http.FileServer(http.Dir("public"))
	http.Handle("/", fs)
	http.HandleFunc("/words", wordHandler)
	fmt.Println("Listening on localhost:8080")
	http.ListenAndServe(":8080", nil)
}
