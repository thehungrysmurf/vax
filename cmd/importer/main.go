package main

import (
	"context"
	"log"

	"github.com/thehungrysmurf/vax/config"
	"github.com/thehungrysmurf/vax/db/store"
	"github.com/thehungrysmurf/vax/importer"

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
		log.Fatalf("failed to connect to database: %v", err)
	}

	dbClient := store.NewDB(conn)

	dataImporter := importer.NewCSVImporter(cfg.VaccinationTotalsFilePath, cfg.ReportsFilePath, cfg.VaccinesFilePath, cfg.SymptomsFilePath, dbClient)
	if err := dataImporter.Run(); err != nil {
		log.Fatalf("failed to importer data: %v", err)
	}

	defer conn.Close(ctx)
}
