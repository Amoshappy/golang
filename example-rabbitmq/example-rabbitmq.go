package main

import (
	"log"
	"net/http"
	"os"

	"github.com/streadway/amqp"
)

// This is a type to hold our word definitions in
// we specifiy both bson (for MongoDB) and json (for web)
// naming for marshalling and unmarshalling
type item struct {
	Word       string `json:"word"`
	Definition string `json:"definition"`
}

// Bind a queue to the exchange to listen for messages
// When we publish a message, it will be sent to this queue, via the exchange
var routingKey = "words"
var exchangeName = "grandtour"
var qName = "sample"
var connection *amqp.Connection
var channel *amqp.Channel
var queue amqp.Queue

func wordHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		log.Println("Getting Message")
		msg, ok, err := channel.Get(qName, true)
		if err != nil {
			log.Fatal(err)
		}

		// Was there something waiting?
		if ok {
			w.Write(msg.Body)
		} else {
			w.Write(nil)
		}
		return
	case "PUT":
		r.ParseForm()
		log.Println("Putting Message")
		err := channel.Publish(exchangeName, routingKey, true, true, amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(r.Form.Get("message")),
		})
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		log.Println("Done")
		w.WriteHeader(http.StatusAccepted)
		return
	}

	return
}

func main() {
	// Connect to database:
	// Connection string in $COMPOSE_RABBITMQ_URL
	var err error

	// amqp.connect(connectionString, { servername: parsedurl.hostname }, function(err, conn) {
	// 	conn.createChannel(function(err, ch) {
	// 	  ch.assertExchange(exchangeName, 'direct', {durable: true});
	// 	  ch.assertQueue(qName, {exclusive: false}, function(err, q) {
	// 		console.log(" [*] Waiting for messages in the queue '%s'", q.queue);
	// 		ch.bindQueue(q.queue, exchangeName, routingKey);
	// 	  });
	// 	});
	// 	setTimeout(function() { conn.close(); }, 500);
	//   });

	log.Println("Dialing")
	connection, err = amqp.Dial(os.Getenv("COMPOSE_RABBITMQ_URL"))
	defer connection.Close()

	if err != nil {
		log.Fatal(err)
	}

	log.Println("Channel")

	channel, err = connection.Channel()

	if err != nil {
		log.Fatal(err)
	}

	log.Println("Exchange Declare")

	err = channel.ExchangeDeclare(exchangeName, "direct", true, false, false, false, nil)

	if err != nil {
		log.Fatal(err)
	}

	log.Println("Queue Declare")

	queue, err = channel.QueueDeclare(qName, true, false, true, false, nil)

	if err != nil {
		log.Fatal(err)
	}

	log.Println("Queue Bind")

	err = channel.QueueBind(queue.Name, routingKey, exchangeName, false, nil)

	if err != nil {
		log.Fatal(err)
	}

	fs := http.FileServer(http.Dir("public"))
	http.Handle("/", fs)
	http.HandleFunc("/message", wordHandler)
	http.ListenAndServe(":8080", nil)
}
