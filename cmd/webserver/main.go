package main

import (
	"context"
	"encoding/json"
	"fmt"
	"internal/cache"
	"log"
	"net/http"
	"sync"
	"time"
)

var (
	c   cache.Cache
	mtx sync.Mutex
)

func main() {
	http.HandleFunc("/setcache", setCacheHandler)
	http.HandleFunc("/set", setHandler)
	http.HandleFunc("/get", getHandler)
	http.HandleFunc("/delete", deleteHandler)
	http.HandleFunc("/description", descriptionHandler)

	log.Println("Starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))

	if c != nil {
		c.Close()
	}
}

func setCacheHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var args struct {
		CacheType string `json:"cacheType"`
		Addr      string `json:"addr"`
		Password  string `json:"password"`
		DB        int    `json:"db"`
	}

	if err := json.NewDecoder(r.Body).Decode(&args); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	newCache, err := cache.NewCache(args.CacheType, args.Addr, args.Password, args.DB)
	if err != nil {
		http.Error(w, "Could not initialize cache: "+err.Error(), http.StatusInternalServerError)
		return
	}

	mtx.Lock()
	c = newCache
	mtx.Unlock()
	fmt.Fprintln(w, "Cache initialized")
}

func setHandler(w http.ResponseWriter, r *http.Request) {
	if c == nil {
		http.Error(w, "Cache not initialized", http.StatusInternalServerError)
		return
	}

	key := r.URL.Query().Get("key")
	value := r.URL.Query().Get("value")
	ttl, _ := time.ParseDuration(r.URL.Query().Get("ttl"))

	if key == "" || value == "" {
		http.Error(w, "Missing key or value", http.StatusBadRequest)
		return
	}

	err := c.Set(context.Background(), key, value, ttl)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "Key %s set with value %s\n", key, value)
}

func getHandler(w http.ResponseWriter, r *http.Request) {
	if c == nil {
		http.Error(w, "Cache not initialized", http.StatusInternalServerError)
		return
	}

	key := r.URL.Query().Get("key")

	if key == "" {
		http.Error(w, "Missing key", http.StatusBadRequest)
		return
	}

	value, err := c.Get(context.Background(), key)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "Key %s has value %v\n", key, value)
}

func deleteHandler(w http.ResponseWriter, r *http.Request) {
	if c == nil {
		http.Error(w, "Cache not initialized", http.StatusInternalServerError)
		return
	}

	key := r.URL.Query().Get("key")

	if key == "" {
		http.Error(w, "Missing key", http.StatusBadRequest)
		return
	}

	err := c.Delete(context.Background(), key)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "Key %s deleted\n", key)
}

func descriptionHandler(w http.ResponseWriter, r *http.Request) {
	if c == nil {
		http.Error(w, "Cache not initialized", http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "%s", c.Description())
}
