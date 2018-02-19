# The Grand Tour - Go

A set of example applications that will add word/definition pairs to a database running on Compose.

Find the articles about the Go section of the Compose Grand Tour here:

* [Go and Compose - MongoDB, Elasticsearch and PostgreSQL](https://www.compose.com/articles/go-and-compose-mongodb-elasticsearch-and-postgresql/)
* [Go and Compose - Redis, RethinkDB, and RabbitMQ](https://www.compose.com/articles/go-and-compose-redis-rethinkdb-and-rabbitmq/)
* [Go and Compose - etcd v3, Scylla, and MySQL](https://www.compose.com/articles/go-and-compose-etcd-v3-scylla-and-mysql/)

This repo contains the apps written in Go. It is intended to run locally.

## Running the Examples

To run from the command-line:
  * navigate to the example-<_db_> directory
  * some examples will just need `go get -a` to be run, others use [glide.sh](http://glide.sh/) or [dep](https://github.com/golang/dep) to vendor libraries where the version needs to be fixed. Consult the example's readme for each example.
  * build the application with `go build`
  * export your Compose connection string as an environment variable 
  * run the application

The application will be served on 127.0.0.1:8080 and can be opened in a browser.

Note: we will be migrating the examples to solely use dep in the future.
