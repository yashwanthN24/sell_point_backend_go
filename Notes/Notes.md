![alt text](image.png)

By default go has its default so we have to set explicitly via the http.Server struct

https://github.com/chickensoft-games/godotenv

To load all environment variables from a .env file into the application we can use the godotenv package

https://pkg.go.dev/github.com/joho/godotenv

internal folder in go is actually a special folder that is used to restrict access to certain packages. Any package placed inside the internal folder can only be imported by other packages that are within the same module or submodule

So if someone tries to import our code as then they wont get access to the interal packages its only within this project

In go there is no concept of private or public packages like in other languages. All packages are public by default, but the internal folder allows us to create a package that is only accessible within the same module or submodule.

In Go there is no centralized repository for packages like npm or pip.
Instead go packages are just the source code that is hosted on a version control system like GitHub, GitLab, or Bitbucket.
When you import a package in Go, you specify the URL of the repository where the package is hosted.

Ex: import "github.com/user/package" (that package name is actually the name of the repository on GitHub)

![alt text](image-1.png)

![alt text](image-2.png)

![alt text](image-3.png)

![alt text](image-4.png)

![alt text](image-5.png)

![alt text](image-6.png)

![alt text](image-7.png)

![alt text](image-8.png)

![alt text](image-9.png)

![alt text](image-10.png)

Must in go funcs is a convention that must always execute the function and always throw a panic when function doesnt execute sucessfully , it is used to load env variables from .env file and if it fails to load it will panic and stop the execution of the program

Here weused MustLoad function (following this convention)

![alt text](image-11.png)

![alt text](image-12.png)

1. pq is a pure Go Postgres driver for the database/sql package. It is used to connect to a PostgreSQL database and execute SQL queries.

![alt text](image-13.png)

2. pgx is a PostgreSQL driver and toolkit for Go. It provides a more feature-rich and flexible interface for working with PostgreSQL databases compared to pq. It supports advanced features like connection pooling, prepared statements, and more. (used in industry modern )

Performance wise pgx is faster than pq because it is designed to be more efficient and optimized for performance. It uses a binary protocol for communication with the PostgreSQL server, which reduces the amount of data that needs to be sent over the network and improves performance.

![alt text](image-14.png)

database/sql is a standard library package in Go that provides a generic interface for interacting with SQL databases. It defines a set of interfaces and types that can be used to execute SQL queries, manage database connections, and handle transactions.

pgx is a PostgreSQL driver that implements the database/sql interfaces, allowing it to be used with the database/sql package. This means that you can use pgx as a drop-in replacement for other database drivers that implement the database/sql interfaces.

https://github.com/jackc/pgx/blob/master/stdlib/sql.go (This we use when we want to use pgx with database/sql interface)

if your sure you want touse only postgres then use

https://github.com/jackc/pgx (This we use when we want to use pgx directly without database/sql interface)

\_ "github.com/jackc/pgx/v5/stdlib" (under import in db.go) and run go mod tidy to install the package

(\_ because we just want to run its init function and dont want to use any of its code )

because we are using pgx as a driver for the existig database/url interface so that in future if we want to change the database we can do it easily without changing the code in db.go file

we just need to change the driver in future code remains same

## Migrations to create tables update schema of table upghrade and dowgrade

we create a separate entry point so running that binary runs the migrations base on whether you call it with `up` or `down` passed as command line arguments

we use `golang-migrate` package for this purpose without using any existsing driver

https://github.com/golang-migrate/migrate#use-in-your-go-project

"github.com/golang-migrate/migrate/v4"
_ "github.com/golang-migrate/migrate/v4/database/postgres"
_ "github.com/golang-migrate/migrate/v4/source/file"

https://github.com/golang-migrate/migrate/blob/master/database/postgres/TUTORIAL.md#optional-run-migrations-within-your-go-app

```go

import (
	"log"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	m, err := migrate.New(
		"file://db/migrations", // pa5th to your migrations folder
		"postgres://postgres:postgres@localhost:5432/example?sslmode=disable" // Databaseconnection string or via env or os.getEnv)
	if err != nil {
		log.Fatal(err)
	}
	if err := m.Up(); err != nil {
		log.Fatal(err)
	}
}

//make migrate-up  (To upgrade the migrate )
// make migrate-down (To downgrade or undo the current version of migration )

```


Still one program m.Down() will undo all migrations up to the first migration so everything is gone 

In production we wont want this behviour we would only want to rollback to last migration so we can use m.Steps(-1) to rollback only one migration at a time

![alt text](image-15.png)

on calling rows.next() it will return true if there is a next row to read and moves to the next row and false if there is no more rows to read

![alt text](image-16.png)

so by calling next on a row we iterate through each row of the query result set and read the data from each row one at a time


## Fetching Database Records & Returning the JSON

Build the listings API endpoint from the previous episode into a fully working, database-backed endpoint. You’ll learn how to query, iterate through query results, scan rows into Go structs, handle database errors safely, and return the final data as properly formatted JSON to the client.


Fetching database data isn't just about running a query. 
A production-ready handler needs to manage the full flow: 

query → iterate → scan → handle errors 
→ release resources → encode → return JSON. 

## Deleting Records & Managing Dependencies
Build the delete listing endpoint and learn how to safely remove a record from the database using a dynamic route parameter and a parameterized SQL query. Then refactor the handlers using Go’s constructor pattern so that dependencies, such as the database, can be injected once and shared across multiple handlers without repeatedly passing them around. 

to get path variables use r.PathValue("pathname")

Ex: r.PathValue("id") (to get the id from the path variable)


## use context to prevent zombie queries 

![alt text](image-17.png)

say db query takes 10s and client disconnects after 5s then the query will still run for 10s and will be a zombie query which is not good for performance

So we have to signal the query to stop executing when the client disconnects so we can use context for that

## Structured Logging 

Structured Logging in Go
Learn how to replace basic text-based logging with structured JSON logging for production Go applications. You’ll use Go’s built-in  log/slog  package, configure a JSON handler, add useful fields such as listing IDs and errors, enable source information, and understand how logging levels control which records are written.
What You’ll Learn
Understand why structured logging matters in production. 

Use Go’s built-in  log/slog  package for structured JSON logging. 

Configure JSON handlers and add useful log attributes. 

Understand  slog  logging levels: Debug, Info, Warn, and Error. 

Configure source information and the default logger. 

Know when to use different logging levels in production.



![alt text](image-18.png)


![alt text](image-19.png)

prints current log level and all other level above it (if we set logginglevel to info then info and all levels above it get printed but debug is ignored even if present in codebase)

![alt text](image-20.png)

![alt text](image-21.png)

loglevel if set to debug then all log messages will be printed, including debug messages. 

This level system works the same across all backends (numbers might change but he same idea)