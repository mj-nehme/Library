[![Makefile CI](https://github.com/mj-nehme/library/actions/workflows/makefile.yml/badge.svg)](https://github.com/mj-nehme/library/actions/workflows/makefile.yml)

# Library CLI

The Library CLI is a command-line tool that allows you to manage books in a system. It provides commands to create, update, delete, list books and perform various other operations.

## A. Installation

To install the Library project, you can use the following steps:

1. **Clone the repository**:

```
git clone git@github.com:mj-nehme/library.git
cd library/
```

2. **Install Prerequisites** (Won't go too much into details here):

- Go 1.20
- make
- Postgres 1.5

3. **Setup Postgres**

Check it online [postgresql website](https://www.postgresql.org/docs/16/tutorial-install.html/)

4. **Setup database**:

4.a. **Create database**:

You can have a look on how to create a database on [postgresql website](https://www.postgresql.org/docs/current/tutorial-createdb.html)

I choosed to name my database: `Library` (You can definitely choose your own name)

4.b. **Create Tables, Relationships and Indices**:


```
psql -U <your_username> -d Library -f ./db/create_db.sql
```

5. **Build**: 

As simple as

```
cd server
make build
```

You can build the API server and API client separately though:

```
make build-server
make build-cli
```

## B. REST API Usage

Now we're all set, let's launch the REST API server:

```
./bin/server
```

Read the `OpenAPI` swagger documentation 