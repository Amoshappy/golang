package main

import (
	"database/sql"
	"encoding/json"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"os"
)

// This is a type to hold our word definitions in
type item struct {
	Word       string `json:"word"`
	Definition string `json:"definition"`
}

var db *sql.DB

func wordHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		rows, err := db.Query("SELECT word,definition FROM words")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		// Create an array of

		var items []*item
		items = make([]*item, 0)
		for rows.Next() {
			myitem := new(item)
			err = rows.Scan(&myitem.Word, &myitem.Definition)
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
	// Connection string in $COMPOSEPOSTGRESQLURL
	// Compse database certificate in composecert.pem
	mydb, err := sql.Open("postgres", os.Getenv("COMPOSEPOSTGRESQLURL"))
	if err != nil {
		log.Fatal(err)
	}
	db = mydb // Copy to global
	defer db.Close()

	_, err = db.Query("CREATE TABLE IF NOT EXISTS words (word varchar(256) NOT NULL, definition varchar(256) NOT NULL)")

	if err != nil {
		log.Fatal(err)
	}

	fs := http.FileServer(http.Dir("public"))
	http.Handle("/", fs)
	http.HandleFunc("/words", wordHandler)
	http.ListenAndServe(":8080", nil)
}
