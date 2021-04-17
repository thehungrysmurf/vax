package main

import (
	"context"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"

	"github.com/thehungrysmurf/vax/config"
	"github.com/thehungrysmurf/vax/db/store"

	"github.com/go-chi/chi/v5"
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

	r := chi.NewRouter()

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("salud"))
	})

	r.Get("/index", func(w http.ResponseWriter, r *http.Request) {
		totals, err := dbClient.GetVaccinationTotals(ctx)
		if err != nil {
			fmt.Fprintf(w, "failed to get vaccination totals %v", err)
		}

		fmt.Fprintf(w, "Totals - Pfizer: %v, Moderna: %v, Janssen: %v", totals.Pfizer, totals.Moderna, totals.Janssen)
	})

	r.Get("/vaccine/{vaccine}", func(w http.ResponseWriter, r *http.Request) {
		var m store.Manufacturer
		vaccine := m.FromString(chi.URLParam(r, "vaccine"))

		counts, err := dbClient.GetSymptomCounts(ctx, vaccine)
		if err != nil {
			fmt.Fprintf(w, "failed to get symptoms %v", err)
		}

		t, err := template.ParseFiles("templates/vaccine.html")
		if err != nil {
			fmt.Fprintf(w, "failed to parse template %v", err)
		}

		ret := VaccinePage{
			Vaccine:       vaccine.String(),
			SymptomCounts: counts,
		}

		if err := t.Execute(w, ret); err != nil {
			fmt.Fprintf(w, "failed to execute template %v", err)
		}
	})

	// TODO return graceful web msg when err != nil in this handler
	r.Get("/vaccine/{vaccine}/category/{name}/{sex}/{agemin}/{agemax}", func(w http.ResponseWriter, r *http.Request) {
		var s store.Sex
		sex := s.FromString(chi.URLParam(r, "sex"))

		ageMin := chi.URLParam(r, "agemin")
		ageFloor, err := strconv.ParseInt(ageMin, 10, 32)
		if err != nil {
			fmt.Fprintf(w, "failed to convert age min to int: %v", err)
		}

		ageMax := chi.URLParam(r, "agemax")
		ageCeil, err := strconv.ParseInt(ageMax, 10, 32)
		if err != nil {
			fmt.Fprintf(w, "failed to convert age min to int: %v", err)
		}

		var m store.Manufacturer
		vaccine := m.FromString(chi.URLParam(r, "vaccine"))

		category := chi.URLParam(r, "name")

		results, err := dbClient.GetFilteredResults(ctx, sex, int(ageFloor), int(ageCeil), vaccine, category)
		if err != nil {
			fmt.Fprintf(w, "failed to get results %v", err)
		}

		t, err := template.ParseFiles("templates/results.html")
		if err != nil {
			fmt.Fprintf(w, "failed to parse template %v", err)
		}

		ret := ResultsPage{
			Vaccine:    vaccine.String(),
			AgeMin:  int(ageFloor),
			AgeMax: int(ageCeil),
			Sex:        sex.String(),
			Results:    results,
		}

		if err := t.Execute(w, ret); err != nil {
			fmt.Fprintf(w, "failed to execute template %v", err)
		}
	})

	log.Fatal(http.ListenAndServe(":8888", r))

	defer conn.Close(ctx)
}

type VaccinePage struct {
	Vaccine string
	SymptomCounts []store.SymptomCount
}

type ResultsPage struct {
	Vaccine string
	AgeMin int
	AgeMax int
	Sex string
	Results []store.FilteredResult
}

type IndexPage struct {
	Pfizer int64
	Moderna int64
	Janssen int64
}
