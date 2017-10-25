package main

import (
	"crypto/tls"
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/garyburd/redigo/redis"
)

// This is a type to hold our word definitions in
// we specify json (for web) naming for marshalling
type item struct {
	Word       string `json:"word"`
	Definition string `json:"definition"`
}

var conn redis.Conn

func wordHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		// Retrieve the Hash Set as a set of key values in a map
		wordresult, err := redis.StringMap(conn.Do("HGETALL", "words"))
		if err != nil {
			log.Fatal(err)
		}
		// Now we need to turn each key value pair into an item so it
		// can be JSON marshalled and sent on
		i := 0
		words := make([]item, len(wordresult))
		for word, def := range wordresult {
			words[i] = item{Word: word, Definition: def}
			i = i + 1
		}

		jsonstr, err := json.Marshal(words)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonstr)
		return
	case "PUT":
		r.ParseForm()
		// Store the new word and definition in a Hash set. Note if the same word
		// is used more than once, it'll overwrite the previous definition.
		_, err := conn.Do("HSET", "words", r.Form.Get("word"), r.Form.Get("definition"))
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
	// Connection string in $COMPOSE_REDIS_URL

	var err error
	uri := os.Getenv("COMPOSE_REDIS_URL")

	// If this is an TLS/SSL connection...
	if strings.HasPrefix(uri, "rediss:") {
		// We need to parse the URI
		parsedURI, err := url.Parse(uri)
		if err != nil {
			log.Fatal(err)
		}
		// To create a TLS servername option because Compose uses SNI
		tlsConfig := &tls.Config{ServerName: parsedURI.Hostname()}
		tlsoption := redis.DialTLSConfig(tlsConfig)
		// And we use that when connecting to let the certificate verify
		conn, err = redis.DialURL(os.Getenv("COMPOSE_REDIS_URL"), tlsoption)
	} else {
		conn, err = redis.DialURL(os.Getenv("COMPOSE_REDIS_URL"))
	}

	if err != nil {
		log.Fatal(err)
	}

	defer conn.Close()

	_, err = redis.String(conn.Do("PING"))

	if err != nil {
		log.Fatal(err)
	}

	fs := http.FileServer(http.Dir("public"))
	http.Handle("/", fs)
	http.HandleFunc("/words", wordHandler)
	http.ListenAndServe(":8080", nil)
}
