# MangaWeb 4 Backend

## Environment

By default, the server will uses environment variables setup inside `.env` and the system for its configurations.

However, you could specify the parameter `--environment=` to use different envionment files. A `dev.env` is supplied in the project, in order to runn the server using this setup, use the command `mangaweb4-backend --environment=dev`

The `dev` environment is setup as follwing:

* it displays log messages in human-readable format, instead of JSON.
* use SQLite as the database.
* read manga files from `./data`

## Setting up Database

Mangaweb 4 supports primarily 2 different kind of databases, SQLite and Postgres. The application does accept MySQL, but it's not extensively tested and thus it's not officialy supported.

By default, Postgres will be used.

To change the database configuration, update two environment variables to suit your database setup. 

**SQLite**

```sh
MANGAWEB_DB_TYPE=sqlite3
MANGAWEB_DB=file:db.sqlite3?cache=shared&_pragma=foreign_keys(1)
```

**Postgres**

```sh
MANGAWEB_DB=postgres
MANGAWEB_DB=postgres://manga:password@host:5432/database
```

Make sure you have setup the database and the user/password beforehand. Also please grant adequate privilege to the user as this user will be used to setup the database (create tables, indexes, etc.). I usually grant all privileges of the database/schema to the user.

## Setup gRPC code generation.

gRPC code is generated from protobuf schema files (*.proto) that is in separated project which is added as a submodule of this project. The code will be generated using `go generate` command. 

However, in contrast to the Ent's geneated code, gRPC requires some elaborated setup in order to get it working.

### Initialize the submodule

To initialize submodule, run the following command.

```sh
$ git submodule init
```

And then, update the submodule to pull the code.

```
$ git submodule update
```

### Install gRPC toolchain

Install `protoc` on the system. On Windows, you could install it using Chocolatey or any similar tool.

```sh
$ choco install protoc
```

Use the following command to install gRPC toolchain for Go.

```sh
$ go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
$ go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
```