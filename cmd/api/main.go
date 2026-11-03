package main

import (
	"fmt"
	"log"
	"net/http"
	_ "net/http/pprof" // this will juts call that package init functions and doesn't give access to objectto acess in this file
	"time"
)

func main() {
	fmt.Println("olx_api running...")

	mux := http.NewServeMux()

	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, r *http.Request) {
		// but we want json apis right
		w.Header().Set("Content-Type", "application/json")

		// w.Write([]byte("all ok")) // by default we are sending text/plain content-type

		w.WriteHeader(http.StatusOK) // this line should comes after all Header set because it sends the header already to the respjnse

		w.Write([]byte(`{"status" : "ok"}`))
	})

	// creating a server with configuration
	srv := http.Server{
		Addr:         ":8090",
		Handler:      mux,
		ReadTimeout:  time.Second * 10,
		WriteTimeout: time.Second * 30,
		IdleTimeout:  time.Second * 60,
	}

	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}

/*


This http.handleFunc Creates a default global router on which asll routes get overidenm so bad practice

create you own router (mux) in go backend always

*/
