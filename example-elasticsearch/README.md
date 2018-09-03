# Compose Grand Tour - Go - Elasticsearch

## Build notes

Before building, run `go get gopkg.in/olivere/elastic.v5` to install the appropriate library.

## Run notes

Set the `COMPOSE_ELASTICSEARCH_URL` environment variable to the Compose connection string for the Elasticsearch database. Remember to create a user for Elasticsearch and include that user's credentials in the URL.

