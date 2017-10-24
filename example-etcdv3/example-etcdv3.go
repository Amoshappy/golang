package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"strings"
	"time"

	"github.com/coreos/etcd/clientv3"
)

type item struct {
	ID         string `json:"_id,omitempty"`
	Word       string `json:"word"`
	Definition string `json:"definition"`
}

var myetcdclient *clientv3.Client

func wordHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		resp, err := myetcdclient.Get(context.TODO(), "/grand_tour/words/", clientv3.WithPrefix())

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
		_, err := myetcdclient.Put(context.TODO(), "/grand_tour/words/"+r.Form.Get("word"), r.Form.Get("definition"))
		if err != nil {
			log.Fatal(err)
		}
		// items = append(items, item{time.Now().Format(time.UnixDate), r.Form.Get("word"), r.Form.Get("definition")})
		w.WriteHeader(http.StatusAccepted)
		return
	}

	return
}

func main() {
	endpointlist := os.Getenv("COMPOSE_ETCD_ENDPOINTS")
	username := os.Getenv("COMPOSE_ETCD_USER")
	password := os.Getenv("COMPOSE_ETCD_PASS")
	certpath := os.Getenv("PATH_TO_ETCD_CERT")

	endpoints := strings.Split(endpointlist, ",")

	var cfg clientv3.Config

	if len(certpath) != 0 {
		// If a certificate path is given, assume a self-signed verification
		// certificate
		// Read the certificate into a file
		caCert, err := ioutil.ReadFile(certpath)
		if err != nil {
			log.Fatal(err)
		}

		// Create a certificate pool
		caCertPool := x509.NewCertPool()
		// and add the freshly read certificate to the new pool
		caCertPool.AppendCertsFromPEM(caCert)

		cfg = clientv3.Config{
			Endpoints:   endpoints,
			TLS:         &tls.Config{RootCAs: caCertPool},
			Username:    username,
			Password:    password,
			DialTimeout: 5 * time.Second,
		}
	} else {
		cfg = clientv3.Config{
			Endpoints:   endpoints,
			Username:    username,
			Password:    password,
			DialTimeout: 5 * time.Second,
		}
	}

	etcdclient, err := clientv3.New(cfg)

	if err != nil {
		log.Fatal(err)
	}

	myetcdclient = etcdclient

	fs := http.FileServer(http.Dir("public"))
	http.Handle("/", fs)
	http.HandleFunc("/words", wordHandler)
	fmt.Println("Listening on localhost:8080")
	http.ListenAndServe(":8080", nil)
}
