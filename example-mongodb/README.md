# Compose Grand Tour - Go - MongoDB

## Build notes

Before building, run `go get gopkg.in/mgo.v2` to install the appropriate library.

## Run notes

Set the `COMPOSE_MONGODB_URL` environment variable to the Compose connection string for the MongoDB database. Remember to create a user in the admin database and include its credentials in the URL.

Set the `PATH_TO_MONGODB_CERT` environment variable to a path to a file containing the Self Signed Certificate
