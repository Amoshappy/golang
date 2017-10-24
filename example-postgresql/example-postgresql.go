package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/lib/pq"
)

// This is a type to hold our word definitions in
type item struct {
	ID         string `json:"id"`
	Word       string `json:"word"`
	Definition string `json:"definition"`
}

var db *sql.DB

func wordHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		rows, err := db.Query("SELECT id,word,definition FROM words")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		items := make([]*item, 0)
		for rows.Next() {
			myitem := new(item)
			err = rows.Scan(&myitem.ID, &myitem.Word, &myitem.Definition)
			if err != nil {
				log.Fatal(err)
			}
			items = append(items, myitem)
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
		_, err := db.Exec("INSERT INTO words (word,definition) VALUES ($1, $2)", r.Form.Get("word"), r.Form.Get("definition"))

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
	// Connection string in $COMPOSE_POSTGRESQL_URL
	// Compose database certificate in $PATH_TO_POSTGRESQL_CERT

	myurl := os.Getenv("COMPOSE_POSTGRESQL_URL") +
		("?sslmode=require&sslrootcert=" + os.Getenv("PATH_TO_POSTGRESQL_CERT"))
	var err error
	db, err = sql.Open("postgres", myurl)
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	_, err = db.Query(`CREATE TABLE IF NOT EXISTS words (
		id serial primary key,
		word varchar(256) NOT NULL, 
		definition varchar(256) NOT NULL)`)

	if err != nil {
		log.Fatal(err)
	}

	fs := http.FileServer(http.Dir("public"))
	http.Handle("/", fs)
	http.HandleFunc("/words", wordHandler)
	fmt.Println("Listening on 8080")
	http.ListenAndServe(":8080", nil)
}
