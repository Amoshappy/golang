package main

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	composeaddresstranslator "github.com/compose/composeaddresstranslator"
	"github.com/gocql/gocql"
	"github.com/scylladb/gocqlx"
	"github.com/scylladb/gocqlx/qb"
)

// This is an item structure we use to hold definitions and to marshal JSON
type item struct {
	ID         gocql.UUID `json:"_id,omitempty" db:"my_table_id"`
	Word       string     `json:"word"`
	Definition string     `json:"definition"`
}

// These are some global variables for the database connection
var cluster *gocql.ClusterConfig
var session *gocql.Session
var addresstranslator composeaddresstranslator.ComposeAddressTranslator

// This is the HTTP request handler for the example. Skip to main() to see the initialisation
// and how this function is configured. Read on to see how data is read and written to the
// database...
func wordHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		stmt, names := qb.Select("examples.words").ToCql()
		q := gocqlx.Query(session.Query(stmt), names)

		var items []item

		if err := gocqlx.Select(&items, q.Query); err != nil {
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
		// First, we parse the incoming form
		r.ParseForm()

		stmt, names := qb.Insert("examples.words").Columns("my_table_id", "word", "definition").ToCql()

		uuid, _ := gocql.RandomUUID()

		newitem := &item{
			uuid,
			r.Form.Get("word"),
			r.Form.Get("definition"),
		}

		q := gocqlx.Query(session.Query(stmt), names).BindStruct(newitem)

		if err := q.ExecRelease(); err != nil {
			log.Fatal(err)
		}

		// And return a good status.
		w.WriteHeader(http.StatusAccepted)
		return
	}
	return
}

// This is the main function where the application starts and sets up
// the database connection, initialises the database and starts the HTTP server

func main() {
	// Get the configuration information from the environment variables
	// Error out if they aren't given.
	urlstring := os.Getenv("COMPOSE_SCYLLA_URLS")
	if len(urlstring) == 0 {
		log.Fatal("No COMPOSE_SCYLLA_URLS given")
	}

	mapstring := os.Getenv("COMPOSE_SCYLLA_MAPS")
	if len(mapstring) == 0 {
		log.Fatal("No COMPOSE_SCYLLA_MAPS given")
	}

	// Extract the first URL to get the username/password combination
	// In this example, thats all that is taken from the URLs; the MAPS provide the rest of the connection data
	urls := strings.Split(urlstring, ",")
	parseurl, err := url.Parse(urls[0])
	if err != nil {
		log.Fatal(err)
	}
	if parseurl == nil {
		log.Fatal("No URL?")
	}
	user := parseurl.User
	username := user.Username()
	password, isset := user.Password()
	if !isset {
		log.Fatal("No Password!")
	}

	// We now use the Compose Address Translator, a specialised version of the address translator which
	// can be initialised with a JSON string.
	addresstranslator, err = composeaddresstranslator.NewFromJSONString(mapstring)
	if err != nil {
		log.Fatal(err)
	}

	// This next command doesn't connect, it just sets up the connection data for connection into a
	// ClusterConfig struct. The first thing it needs though is the addresses of the cluster. Handily,
	// the ComposeAddressTranslator does the work for you, turning the map into ContactPoints so...
	cluster = gocql.NewCluster(addresstranslator.ContactPoints()...)
	// You might expect a cluster.Keyspace = "example"  here but... If we specified a default keyspace
	// and it didn't exist, the example would error out before trying to create it. So we don't.
	cluster.Consistency = gocql.Quorum
	cluster.SslOpts = &gocql.SslOptions{} // Turns on SSL
	cluster.Authenticator = gocql.PasswordAuthenticator{
		Username: username,
		Password: password,
	}
	cluster.AddressTranslator = addresstranslator
	cluster.IgnorePeerAddr = true

	// Now we create a session and it's here that we connect to the database
	session, err = cluster.CreateSession()
	if err != nil {
		log.Fatal(err)
	}
	// We also set up a defer to close the session when done.
	defer session.Close()

	// Now we create a new keyspace if we need it
	err = session.Query("CREATE KEYSPACE IF NOT EXISTS examples WITH replication = {'class': 'SimpleStrategy', 'replication_factor': '3' }").Exec()
	if err != nil {
		log.Fatal(err)
	}

	// And similarly we create the table if we need it too
	err = session.Query("CREATE TABLE IF NOT EXISTS examples.words (my_table_id uuid, word text, definition text, PRIMARY KEY(my_table_id))").Exec()
	if err != nil {
		log.Fatal(err)
	}

	// Finally, for main, we set up the HTTP server. Incoming requests will be dealt with
	// by the wordHandler function
	fs := http.FileServer(http.Dir("public"))
	http.Handle("/", fs)
	http.HandleFunc("/words", wordHandler)
	http.ListenAndServe(":8080", nil)
}
