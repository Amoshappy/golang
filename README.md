# The Grand Tour - Go

A set of example applications that will add word/definition pairs to a database running on Compose.

This repo contains the apps written in Go. It is intended to run locally, for example and experimentation purposes only.

## Running the Example
To run from the command-line:
  * navigate to the example-<_db_> directory
  * some examples will just need `go get -a` to be run, others use [glide.sh](http://glide.sh/) to vendor libraries where the version needs to be fixed. Consult the readme for each example.
  * build the application with `go build`
  * export your Compose connection string as an environment variable 
  * run the application

The application will be served on 127.0.0.1:8080 and can be opened in a browser.
