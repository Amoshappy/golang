package main

import (
	"context"
	"encoding/json"
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

var myclient *elastic.Client

func wordHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		searchResult, err := myclient.Search().
			Index("example").
			Type("word").
			//		Sort("added", false).
			Do(context.TODO())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		words := []item{}
		var i item
		for _, peritem := range searchResult.Each(reflect.TypeOf(i)) {
			i := peritem.(item)
			words = append(words, i)
		}
		newjson, err := json.Marshal(words)
		w.Write(newjson)

		return
	case "PUT":
		r.ParseForm()
		newitem := item{r.Form.Get("word"),
			r.Form.Get("definition")}

		_, err := myclient.Index().
			Index("example").
			Type("word").
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
	esuri := os.Getenv("COMPOSEESURI")

	client, err := elastic.NewClient(elastic.SetURL(esuri), elastic.SetSniff(false))

	myclient = client

	if err != nil {
		log.Fatal(err)
	}

	exists, err := client.IndexExists("examples").Do(context.TODO())
	if err != nil {
		log.Fatal(err)
	}

	if !exists {
		_, err := client.CreateIndex("examples").Do(context.TODO())
		if err != nil {
			log.Fatal(err)
		}
	}

	fs := http.FileServer(http.Dir("public"))
	http.Handle("/", fs)
	http.HandleFunc("/words", wordHandler)
	http.ListenAndServe(":8080", nil)
}
