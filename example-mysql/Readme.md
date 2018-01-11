# Compose Grand Tour Go/MySQL Example

Two Grand Tour examples of connecting to MySQL and reading/writing data are included here.

The first, `example-mysql`, is the definitive example using only the golang SQL and mysql driver.

The second, `example-mysql-sqlx`, uses the sqlx extension to golang sql to make connected and reading of the data more elegant.

## Building

This example comes with vendored dependencies, managed by [`dep`](https://github.com/golang/dep). Install that and run `dep ensure` to check and install any missing dependencies.

To build each application run either

`go build example-mysql.go`

or

`go build example-mysql-sqlx.go`

## Before running

One environment variables must be set: `COMPOSE_MYSQL_URL`.

* COMPOSE_MYSQL_URL - the Compose connection string for the MySQL database. Remember to create a user for ScyllaDB and include that user's credentials in the URL.

### Examples

```
export COMPOSE_MYSQL_URL="mysql://admin:password@sl-eu-lon-2-portal.0.dblayer.com:17851/compose"
```

## Running

Once set, run either `./example-mysql` or `./example-scylla-mysql-sqlx` and point your browser at localhost:8080.

