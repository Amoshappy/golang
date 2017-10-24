package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	elastic "gopkg.in/olivere/elastic.v5"

	"reflect"
)

type item struct {
	Word       string `json:"word"`
	Definition string `json:"definition"`
}

var client *elastic.Client

func wordHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		searchResult, err := client.Search().
			Index("grand_tour").
			Type("words").
			Do(context.TODO())

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		words := []item{}
		var i item
		for _, peritem := range searchResult.Each(reflect.TypeOf(i)) {
			i := peritem.(item)
			words = append(words, i)
		}

		w.Header().Set("Content-Type", "application/json")
		newjson, err := json.Marshal(words)
		w.Write(newjson)

		return
	case "PUT":
		r.ParseForm()
		newitem := item{r.Form.Get("word"),
			r.Form.Get("definition")}

		_, err := client.Index().
			Index("grand_tour").
			Type("words").
			BodyJson(newitem).
			Refresh("true").
			Do(context.TODO())
		if err != nil {
			log.Fatal(err)
		}
		w.WriteHeader(http.StatusAccepted)
		return
	}

	return
}

func main() {
	esuri := os.Getenv("COMPOSE_ELASTICSEARCH_URL")

	var err error
	client, err = elastic.NewClient(elastic.SetURL(esuri), elastic.SetSniff(false))

	if err != nil {
		log.Fatal(err)
	}

	// Check if index exists.
	exists, err := client.IndexExists("grand_tour").Do(context.TODO())
	if err != nil {
		log.Fatal(err)
	}

	// If not, create it
	if !exists {
		_, err := client.CreateIndex("grand_tour").Do(context.TODO())
		if err != nil {
			log.Fatal(err)
		}
	}

	fs := http.FileServer(http.Dir("public"))
	http.Handle("/", fs)
	http.HandleFunc("/words", wordHandler)
	fmt.Println("Listening on localhost:8080")
	http.ListenAndServe(":8080", nil)
}
