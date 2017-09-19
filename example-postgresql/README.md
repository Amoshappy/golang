# Compose Grand Tour - Go - PostgreSQL

## Build notes

Before building, run `go get -a` to install the lib/pq library.
Build with `go build`.


## Run notes

Set the `COMPOSE_POSTGRESQL_URL` environment variable to the Compose connection string for the PostgreSQL database. 
Set the `PATH_TO_POSTGRESQL_CERT` to the full path of the self-signed certificate for the deployment.


