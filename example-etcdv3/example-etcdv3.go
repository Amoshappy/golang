package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path"
	"strings"
	"time"

	"github.com/coreos/etcd/clientv3"
)

type item struct {
	Word       string `json:"word"`
	Definition string `json:"definition"`
}

var etcdclient *clientv3.Client

func wordHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		resp, err := etcdclient.Get(context.TODO(), "/grand_tour/words/", clientv3.WithPrefix())

		if err != nil {
			log.Fatal(err)
		}

		var items []item

		for _, ev := range resp.Kvs {
			_, word := path.Split(string(ev.Key))
			items = append(items, item{Word: word, Definition: string(ev.Value)})
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
		_, err := etcdclient.Put(context.TODO(), "/grand_tour/words/"+r.Form.Get("word"), r.Form.Get("definition"))
		if err != nil {
			log.Fatal(err)
		}
		w.WriteHeader(http.StatusAccepted)
		return
	}

	return
}

func main() {
	var err error

	endpointlist := os.Getenv("COMPOSE_ETCD_ENDPOINTS")
	username := os.Getenv("COMPOSE_ETCD_USER")
	password := os.Getenv("COMPOSE_ETCD_PASS")

	endpoints := strings.Split(endpointlist, ",")

	cfg := clientv3.Config{
		Endpoints:   endpoints,
		Username:    username,
		Password:    password,
		DialTimeout: 5 * time.Second,
	}

	etcdclient, err = clientv3.New(cfg)

	if err != nil {
		log.Fatal(err)
	}

	fs := http.FileServer(http.Dir("public"))
	http.Handle("/", fs)
	http.HandleFunc("/words", wordHandler)
	fmt.Println("Listening on localhost:8080")
	http.ListenAndServe(":8080", nil)
}
