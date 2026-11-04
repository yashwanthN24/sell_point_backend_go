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
