package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

type item struct {
	ID         string `json:"_id,omitempty"`
	Word       string `json:"word"`
	Definition string `json:"definition"`
}

var db *sql.DB

func wordHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		// Query the database for the rows
		rows, err := db.Query("SELECT id,word,definition FROM words ORDER BY word ASC")
		if err != nil {
			log.Fatal(err)
		}
		defer rows.Close()
		// Now create somewhere to hold the rows
		items := make([]*item, 0)
		// Read each row
		for rows.Next() {
			// Make an empty item
			myitem := new(item)
			// Scan results into that item
			err := rows.Scan(&myitem.ID, &myitem.Word, &myitem.Definition)
			if err != nil {
				log.Fatal(err)
			}
			// Then append it to our array
			items = append(items, myitem)
		}
		// Make sure we finished cleanly
		err = rows.Err()
		if err != nil {
			log.Fatal(err)
		}
		// Now convert the array into JSON and send it as the response.
		jsonstr, err := json.Marshal(items)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonstr)
		return
	case "PUT":
		// First, we parse the incoming form
		r.ParseForm()
		// Now we simply INSERT the data into the table
		// The id generates itself automatically
		_, err := db.Exec("INSERT INTO words(word, definition) VALUES(?,?)", r.Form.Get("word"), r.Form.Get("definition"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		// And return a good status.
		w.WriteHeader(http.StatusAccepted)
		return
	}

	return
}

func main() {
	connectionString := os.Getenv("COMPOSE_MYSQL_URL")

	url, err := url.Parse(connectionString)
	if err != nil {
		log.Fatal(err)
	}

	passwd, _ := url.User.Password()

	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s",
		url.User.Username(),
		passwd,
		url.Host,
		url.Path[1:])

	db, err = sql.Open("mysql", dsn)

	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec("CREATE TABLE IF NOT EXISTS words (id int auto_increment primary key, word varchar(256) NOT NULL, definition varchar(256) NOT NULL)")
	if err != nil {
		log.Fatal(err)
	}

	fs := http.FileServer(http.Dir("public"))
	http.Handle("/", fs)
	http.HandleFunc("/words", wordHandler)
	http.ListenAndServe(":8080", nil)
}
