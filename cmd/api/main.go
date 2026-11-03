package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	// dotenv "github.com/joho/godotenv" // normal available as godotenv but now we can name as them dotenv so now all avilable via dotenv

	"github.com/YashwanthN24/sell_point/internal/config"
)

func main() {

	cfg := config.MustLoad()
	fmt.Println("starting olx server...")

	mux := http.NewServeMux()

	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, r *http.Request) {
		// but we want json apis right
		w.Header().Set("Content-Type", "application/json")

		// w.Write([]byte("all ok")) // by default we are sending text/plain content-type

		w.WriteHeader(http.StatusOK) // this line should comes after all Header set because it sends the header already to the respjnse

		w.Write([]byte(`{"status" : "ok"}`))
	})

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
