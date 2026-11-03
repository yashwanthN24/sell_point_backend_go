package main

import (
	"fmt"
	"log"
	"net/http"
	_ "net/http/pprof" // this will juts call that package init functions and doesn't give access to objectto acess in this file
)

func main() {
	fmt.Println("olx_api running...")

	http.HandleFunc("GET /healthz", func(w http.ResponseWriter, r *http.Request) {
		// but we want json apis right
		w.Header().Set("Content-Type", "application/json")

		// w.Write([]byte("all ok")) // by default we are sending text/plain content-type

		w.WriteHeader(http.StatusOK)

		w.Write([]byte(`{"status" : "ok"}`))
	})

	err := http.ListenAndServe(":8090", nil)
	if err != nil {
		log.Fatalf("server failed: %v", err)
	}
}

/*


This http.handleFunc Creates a default global router on which asll routes get overidenm so bad practice 

create you own router (mux) in go backend always 

*/