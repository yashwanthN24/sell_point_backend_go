package main

import (
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"time"

	// dotenv "github.com/joho/godotenv" // normal available as godotenv but now we can name as them dotenv so now all avilable via dotenv

	"github.com/YashwanthN24/sell_point/internal/config"
	"github.com/YashwanthN24/sell_point/internal/db"
	"github.com/YashwanthN24/sell_point/internal/handlers"
	"github.com/YashwanthN24/sell_point/internal/middlewares"
)

func main() {

	cfg := config.MustLoad()
	db, err := db.Connect(cfg.DatabaseUrl)

	// using this db variable we can run queries

	if err != nil {
		log.Fatalf("main.db.connect : %v", err)
	}

	logHandler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,           // log even the file and speific exact line where this log is from
		Level:     slog.LevelInfo, // set global handler so global olog levell is Info (so info , warn and debug logs are logged rest al ignore i.e debug logs avoided )
	})
	logger := slog.New(logHandler)
	slog.SetDefault(logger) // first this sets slog default handler as JsonLogger and this make sure that default log even log will be redirected to slog json style logger
	fmt.Println("database connected")
	fmt.Println("starting olx server...")

	mux := http.NewServeMux()

	// refactor handlers depednednyc injection via constructor pattern using a common struct that includes all the dependency  of that handler
	lh := handlers.NewListingHandler(db, logger)
	ah := handlers.NewAuthHandler(db, logger , cfg)

	mux.HandleFunc("GET /healthz", handlers.Health)
	mux.HandleFunc("GET /listings", lh.List)

	mux.HandleFunc("DELETE /listings/{id}", lh.Delete)

	mux.HandleFunc("POST /listings", lh.Create)

	mux.HandleFunc("POST /signup", ah.Signup)

	mux.HandleFunc("POST /signin" , ah.Signin)

	// so all routes go through this middleware before going to the actual route handler
	handler := middlewares.RequestId(mux)
	// creating a server with configuration
	srv := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      handler,
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
