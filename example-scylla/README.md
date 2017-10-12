# Compose Grand Tour Go/Scylla Example

Two Grand Tour examples of connecting to Scylla and reading/writing data are included here.

The first, `example-scylla`, is the definitive example using only the gocql driver.

The second, `example-scylla-gocqlx`, uses the gocqlx extension to gocql to make the reading and writing of data more elegant.

## Building

This example comes with vendored dependencies, managed by `[dep])(https://github.com/golang/dep)`. Install that and run `dep ensure` to check and install any missing dependencies.

To build each application run either

`go build example-scylla.go`

or

`go builde example-scylla-gocqlx.go`

## Running

Two environment variables must be set: `COMPOSE_SCYLLA_URLS` and `COMPOSE_SCYLLA_MAPS`

* COMPOSE_SCYLLADB_URL - the Compose connection string for the ScyllaDB database. Remember to create a user for ScyllaDB and include that user's credentials in the URL.
* COMPOSE_SCYLLADB_MAPS - the Address Translation Map for the ScyllaDB database. Copy the full contents as shown on your deployment's overview page.

Examples

```
export COMPOSE_SCYLLA_URLS="https://scylla:password@portal1122-5.regal-scylla-68.compose-3.composedb.com:20598,https://scylla:password@portal1085-4.regal-scylla-68.compose-3.composedb.com:20598,https://scylla:password@portal1130-0.regal-scylla-68.compose-3.composedb.com:20598"
export COMPOSE_SCYLLA_MAPS='{
                  "10.153.168.133:9042": "portal1122-5.regal-scylla-68.compose-3.composedb.com:20598",
                  "10.153.168.134:9042": "portal1085-4.regal-scylla-68.compose-3.composedb.com:20598",
                  "10.153.168.135:9042": "portal1130-0.regal-scylla-68.compose-3.composedb.com:20598"
                }'
```

Once set, run either `./example-scylla` or `./example-scylla-gocqlx` and point your browser at localhost:8080.
