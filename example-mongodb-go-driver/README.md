# Compose Grand Tour - Go - MongoDB

This Grand Tour MongoDB/Go is a preview release using an alpha release of the [MongoDB Go Driver](https://github.com/mongodb/mongo-go-driver). Expect it and the driver to change, but this should give an early impression of the future of MongoDB and Go.

## Build notes

Before building, either run `go get github.com/mongodb/mongo-go-driver/mongo` to install the appropriate package or use `dep` and run `dep ensure -add github.com/mongodb/mongo-go-driver/mongo` to add the dependency.

## Run notes

Set the `COMPOSE_MONGODB_URL` environment variable to the Compose connection string for the MongoDB database. Remember to create a user in the admin database and include its credentials in the URL.

Set the `PATH_TO_MONGODB_CERT` environment variable to a path to a file containing the Self Signed Certificate
