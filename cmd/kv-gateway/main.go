package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/DevPatel-11/distributed-raft-kv-store/kvstore"
)

var store *kvstore.Store

func main() {
	gatewayAddr := flag.String("addr", "0.0.0.0:8080", "Gateway address")
	flag.Parse()

	fmt.Printf("Starting KV Gateway on %s...\n", *gatewayAddr)

	// Initialize the store
	store = kvstore.New()

	// Register HTTP handlers
	http.HandleFunc("/kv/", handleKV)
	http.HandleFunc("/health", handleHealth)

	// Start the HTTP server
	log.Printf("Gateway listening on %s", *gatewayAddr)
	if err := http.ListenAndServe(*gatewayAddr, nil); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}

func handleKV(w http.ResponseWriter, r *http.Request) {
	// Extract key from URL path: /kv/{key}
	key := strings.TrimPrefix(r.URL.Path, "/kv/")
	if key == "" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Key is required")
		return
	}

	w.Header().Set("Content-Type", "application/json")

	switch r.Method {
	case http.MethodGet:
		handleGet(w, key)
	case http.MethodPut:
		handlePut(w, r, key)
	case http.MethodDelete:
		handleDelete(w, key)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprintf(w, "Method not allowed")
	}
}

func handleGet(w http.ResponseWriter, key string) {
	value, ok := store.Get(key)
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "key not found"})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"key": key, "value": string(value)})
}

func handlePut(w http.ResponseWriter, r *http.Request, key string) {
	var req map[string]string
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "invalid request"})
		return
	}

	value, ok := req["value"]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "value is required"})
		return
	}

	if err := store.Set(key, []byte(value)); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"key": key, "status": "created"})
}

func handleDelete(w http.ResponseWriter, key string) {
	if err := store.Delete(key); err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"key": key, "status": "deleted"})
}

func handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "healthy", "keys": fmt.Sprintf("%d", store.Size())})
}
