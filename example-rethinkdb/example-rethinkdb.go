package main

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"

	rethink "gopkg.in/gorethink/gorethink.v2"
)

// This is a type to hold our word definitions in
type item struct {
	ID         string `gorethink:"id,omitempty" json:"_id,omitempty"`
	Word       string `gorethink:"word" json:"word"`
	Definition string `gorethink:"definition" json:"definition"`
}

var session *rethink.Session

func wordHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		results, err := rethink.DB("examples").Table("words").Run(session)

		var items []*item
		err = results.All(&items)
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
		newitem := item{Word: r.Form.Get("word"), Definition: r.Form.Get("definition")}

		err := rethink.DB("examples").Table("words").Insert(newitem).Exec(session)

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
	// Connection string in $COMPOSE_RETHINKDB_URL
	// Compse database certificate in $PATH_TO_RETHINKDB_CERT
	roots := x509.NewCertPool()
	cert, err := ioutil.ReadFile(os.Getenv("PATH_TO_RETHINKDB_CERT"))
	if err != nil {
		log.Fatal(err)
	}
	roots.AppendCertsFromPEM(cert)

	rethinkurl, err := url.Parse(os.Getenv("COMPOSE_RETHINKDB_URL"))

	if err != nil {
		log.Fatal(err)
	}

	password, setpass := rethinkurl.User.Password()

	if !setpass {
		log.Fatal("Password needs to be set in $COMPOSE_RETHINKDB_URL")
	}

	session, err = rethink.Connect(rethink.ConnectOpts{
		Address:  rethinkurl.Host,
		Username: rethinkurl.User.Username(),
		Password: password,
		TLSConfig: &tls.Config{
			RootCAs: roots,
		},
	})

	if err != nil {
		log.Fatalln(err)
	}

	rethink.DBCreate("examples").Exec(session)
	rethink.DB("examples").TableCreate("words").Exec(session)
	defer session.Close()

	fs := http.FileServer(http.Dir("public"))
	http.Handle("/", fs)
	http.HandleFunc("/words", wordHandler)
	http.ListenAndServe(":8080", nil)
}
