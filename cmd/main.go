package main

import (
	"context"
	"log"

	"github.com/thehungrysmurf/vax/config"
	"github.com/thehungrysmurf/vax/data"

	"github.com/jackc/pgx/v4"
	"github.com/joeshaw/envdecode"
)

func main() {
	var cfg config.Config
	err := envdecode.Decode(&cfg)
	if err != nil {
		log.Fatalf("failed to read config: %v", err)
	}

	dataImporter := data.NewCSVImporter(cfg.ReportsFilePath, cfg.VaccinesFilePath, cfg.SymptomsFilePath)
	if err := dataImporter.Run(); err != nil {
		log.Fatalf("failed to import data: %v", err)
	}

	ctx := context.Background()
	conn, err := pgx.Connect(ctx, cfg.DatabaseURI)
	if err != nil {
		log.Fatal("failed to connect to database")
	}

	defer conn.Close(ctx)
}
