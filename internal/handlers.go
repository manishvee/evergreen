package internal

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
)


type CreateIndexRequest struct {
	IndexName string `json:"name"`
}

func CreateIndexHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req CreateIndexRequest

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	err := NewFileStore(req.IndexName)
	if errors.Is(err, ErrIndexAlreadyExists) {
		http.Error(w, "Index already exists", http.StatusConflict)
		return
	} else if err != nil {
		msg := "Failed to create index"
		slog.Error(msg, "err", err)
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"status": "index created"})
}

func RootHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "test")
}
