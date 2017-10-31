package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/streadway/amqp"
)

// Bind a queue to the exchange to listen for messages
// When we publish a message, it will be sent to this queue, via the exchange
var routingKey = "words"
var exchangeName = "grandtour"
var connection *amqp.Connection
var channel *amqp.Channel
var queue amqp.Queue

func wordHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		// Use synchronous Get to get a single message from
		// the queue. In real life, prefer to Consume.
		msg, ok, err := channel.Get(queue.Name, true)
		if err != nil {
			log.Fatal(err)
		}

		// Was there something waiting?
		if ok {
			// Write the retrieved payload straight into the body
			w.Write(msg.Body)
		} else {
			// Or if nothing was retrieved, a message reflecting that
			w.Write([]byte("No message waiting"))
		}
		return
	case "PUT":
		r.ParseForm()
		// Compose a message comprising of the form message and the time
		// now
		msg := r.Form.Get("message") + " " + time.Now().Format("15:04:05.00")

		// Publish the messahe to the exchange with our routing key
		err := channel.Publish(exchangeName,
			routingKey,
			false,
			false,
			amqp.Publishing{
				ContentType: "text/plain",
				Body:        []byte(msg),
			})

		if err != nil {
			log.Fatal(err)
		}
		// Set the response to accepted
		w.WriteHeader(http.StatusAccepted)
		// and rite the message that was sent back to the user
		w.Write([]byte(msg))
		return
	}

	return
}

func main() {
	// Connect to database:
	// Connection string in $COMPOSE_RABBITMQ_URL
	var err error

	// The library handles amqps connections and sets Servername correctly
	// so we can just connect here with the Dial() function
	connection, err = amqp.Dial(os.Getenv("COMPOSE_RABBITMQ_URL"))
	if err != nil {
		log.Fatal(err)
	}
	defer connection.Close()

	// We create a channel to communicate over. Note, if there is an
	// error on the channel, we would need to recreate the channel.
	channel, err = connection.Channel()
	if err != nil {
		log.Fatal(err)
	}

	// We can now create a durable exchange where we'll post our messages
	err = channel.ExchangeDeclare(exchangeName, "direct", true, false, false, false, nil)
	if err != nil {
		log.Fatal(err)
	}

	// And we'll declare a transient queue for our client
	queue, err = channel.QueueDeclare("", false, false, false, false, nil)
	if err != nil {
		log.Fatal(err)
	}

	// Now we'll bind out transient queue to recieve messages with our
	// routing key
	err = channel.QueueBind(queue.Name, routingKey, exchangeName, false, nil)
	if err != nil {
		log.Fatal(err)
	}

	// With setup done we begin web serving
	fs := http.FileServer(http.Dir("public"))
	http.Handle("/", fs)
	http.HandleFunc("/message", wordHandler)
	http.ListenAndServe(":8080", nil)
}
