package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	// dotenv "github.com/joho/godotenv" // normal available as godotenv but now we can name as them dotenv so now all avilable via dotenv

	"github.com/YashwanthN24/sell_point/internal/config"
	"github.com/YashwanthN24/sell_point/internal/db"
	"github.com/YashwanthN24/sell_point/internal/handlers"
)

func main() {

	cfg := config.MustLoad()
	_, err := db.Connect(cfg.DatabaseUrl)

	// using this db variable we can run queries 

	if err != nil {
		log.Fatalf("main.db.connect : %v", err)
	}



	fmt.Println("database connected")
	fmt.Println("starting olx server...")

	mux := http.NewServeMux()

	mux.HandleFunc("GET /healthz", handlers.Health)

	// creating a server with configuration
	srv := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      mux,
		ReadTimeout:  time.Second * 10,
		WriteTimeout: time.Second * 30,
		IdleTimeout:  time.Second * 60,
	}

	log.Printf("server is listening on %s", srv.Addr)
	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}

/*


This http.handleFunc Creates a default global router on which asll routes get overidenm so bad practice

create you own router (mux) in go backend always

*/
