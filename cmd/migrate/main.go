package main

import (
	"fmt"
	"log"
	"os"

	"github.com/YashwanthN24/sell_point/internal/config"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("usage: migrate <up | down>")
	}

	cfg := config.MustLoad()

	m, err := migrate.New(
		"file://migrations", // path to your  migration files
		cfg.DatabaseUrl,     // DatabaseUrl
	)

	if err != nil {
		log.Fatalf("migration.New : %v ", err)
	}

	switch os.Args[1] {
	case "up":
		log.Println("up called")
		if err := m.Up(); err != nil {
			log.Fatal(err)
		}
	case "down":
		log.Println("down called")
		if err := m.Steps(-1); err != nil { // downgrade by 1 migrations version not entire whole migration till the start which can be dangerous as it removes everyy m,igration from the start of application
			log.Fatal(err)
		}
	default:
		log.Fatalf("unknown command : %s", os.Args[1])
	}

	fmt.Println("Running migrations...")
}
