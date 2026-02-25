package main

import (
	"log"
	"log/slog"
	"net/http"

	"github.com/manishvee/evergreen/internal"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", internal.RootHandler)
	mux.HandleFunc("POST /indexes", internal.CreateIndexHandler)

	slog.Info("server listening on port 5225")
	log.Fatal(http.ListenAndServe(":5225", mux))
}
