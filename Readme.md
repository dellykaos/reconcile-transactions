# Reconciliation Service

This service is used to identify unmatched and discrepant transactions between internal data (system transactions) and external data (bank statements).

The tech stack used for this services are:

- Golang version 1.23.0
- PostgreSQL 14.8

## How to run

### 1. Setup environment

Copy file `env.sample` to `.env`

```shell
cp env.sample .env
```

Fill all required configuration needed.

### 2. Setup database

Please make sure you have PostgreSQL 14.8 installed and `golang-migrate` CLI for PostgreSQL is installed.

To Install `golang-migrate` CLI:

```shell
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@v4.17.1
```

(Optional)
If you want to use dockerized PostgreSQL, you can setup it by using these command.

```shell
docker pull postgres:14.8-alpine3.18 # pull docker image
docker run --name postgres14 -d -p 5432:5432 -e POSTGRES_PASSWORD=your_password -e POSTGRES_USER=your_user postgres:14.8-alpine3.18 # run docker image postgres14
```

If the database is already started, run this to setup the database:

```shell
make setup-db
```

To run the DB Migration individually:

```shell
make migrate-up
```

If you need to rollback the database, you can use this command:

```shell
make migrate-down
```
