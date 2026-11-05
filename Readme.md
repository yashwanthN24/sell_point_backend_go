## Go Project Setup

- go mod init `<project_name>`

Convention

- go mod init `github.com/<username>/<project_name></project_name>`

- main.go entry point for the application contains the main function.

- go run main.go (builds and runs the application in one go , good for development)

- For productionwe wnt to build the application and run the binary.

- `go build  -o <binary_name> main.go` (builds the application and creates a binary with the name `<binary_name>`)

- ./<binary_name> (runs the binary)

to automate this create a Makefile with the following content

```Makefile
.PHONY: build run

build:
    @go build -o bin/main main.go

run: build
    @./bin/main

```

@ to avoid that command being shown in the terminal

.PHONY to tell that stages are commands not file/folder names meanung build and run are make commads not folder/file names

make run (runs even build because of run: build line in Makefile that tells make to run build stage first before running run)

make build

Right now it compiles all go files at the root file including the main.go fiel to one binary 

but if we have multiple entry points for our app say ws , cli or worker 

then we have to create cmd folder and create indicual foldersfor each entry point and have its own main.go file having main function for each folder
