package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/thehungrysmurf/vax/config"
	"github.com/thehungrysmurf/vax/db/store"

	"github.com/jackc/pgx/v4"
	"github.com/joeshaw/envdecode"
)

func main() {
	var cfg config.Config
	err := envdecode.Decode(&cfg)
	if err != nil {
		log.Fatalf("failed to read config: %v", err)
	}

	ctx := context.Background()
	conn, err := pgx.Connect(ctx, cfg.DatabaseURI)
	if err != nil {
		log.Fatal("failed to connect to database")
	}

	dbClient := store.NewDB(conn)

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "OK")
	})

	http.HandleFunc("/counts", func(w http.ResponseWriter, r *http.Request) {
		counts, err := dbClient.GetSymptomCounts(ctx, store.Moderna)
		if err != nil {
			fmt.Fprintf(w, "Got an error %v", err)
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(counts); err != nil {
			fmt.Fprintf(w, "Got an error encoding json to return %v", err)
		}
	})

	http.HandleFunc("/results", func(w http.ResponseWriter, r *http.Request) {
		results, err := dbClient.GetFilteredResults(ctx, store.Female, 12, 99, store.Moderna, "flu-like")
		if err != nil {
			fmt.Fprintf(w, "Got an error %v", err)
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(results); err != nil {
			fmt.Fprintf(w, "Got an error encoding json to return %v", err)
		}
	})

	log.Fatal(http.ListenAndServe(":8888", nil))

	defer conn.Close(ctx)
}
