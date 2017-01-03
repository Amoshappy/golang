package main

import (
	"encoding/json"
	"github.com/garyburd/redigo/redis"
	"log"
	"net/http"
	"os"
)

// This is a type to hold our word definitions in
// we specifiy both bson (for MongoDB) and json (for web)
// naming for marshalling and unmarshalling
type item struct {
	Word       string `json:"word"`
	Definition string `json:"definition"`
}

var conn redis.Conn

func wordHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		wordresult, err := redis.StringMap(conn.Do("HGETALL", "words"))
		// This is a map so we need to make it into an array of items
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
	// Connection string in $COMPOSEREDISBURL
	var err error

	conn, err = redis.DialURL(os.Getenv("COMPOSEREDISURL"))

	if err != nil {
		log.Fatal(err)
	}

	fs := http.FileServer(http.Dir("public"))
	http.Handle("/", fs)
	http.HandleFunc("/words", wordHandler)
	http.ListenAndServe(":8080", nil)
}
