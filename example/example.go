package main

import (
	"encoding/json"
	"net/http"
	"time"
)

type item struct {
	ID         string `json:"_id,omitempty"`
	Word       string `json:"word"`
	Definition string `json:"definition"`
}

var items []item

func wordHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
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
		items = append(items, item{time.Now().Format(time.UnixDate), r.Form.Get("word"), r.Form.Get("definition")})
		w.WriteHeader(http.StatusAccepted)
		return
	}

	return
}

func main() {
	items = append(items, item{"1", "hello", "a greeting"})

	fs := http.FileServer(http.Dir("public"))
	http.Handle("/", fs)
	http.HandleFunc("/words", wordHandler)
	http.ListenAndServe(":8080", nil)
}
