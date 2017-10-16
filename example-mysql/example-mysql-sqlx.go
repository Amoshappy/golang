package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

type item struct {
	ID         string `json:"_id,omitempty"`
	Word       string `json:"word"`
	Definition string `json:"definition"`
}

var db *sqlx.DB

func wordHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		items := []item{}
		err := db.Select(&items, "SELECT * FROM words ORDER BY word ASC")
		if err != nil {
			log.Fatal(err)
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
		_, err := db.Exec("INSERT INTO words(word, definition) VALUES(?,?)", r.Form.Get("word"), r.Form.Get("definition"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
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

	db, err = sqlx.Connect("mysql", dsn)

	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	_, err = db.Exec("CREATE TABLE IF NOT EXISTS words (id int auto_increment primary key, word varchar(256) NOT NULL, definition varchar(256) NOT NULL)")
	if err != nil {
		log.Fatal(err)
	}

	fs := http.FileServer(http.Dir("public"))
	http.Handle("/", fs)
	http.HandleFunc("/words", wordHandler)
	http.ListenAndServe(":8080", nil)
}
