package main

import (
	"context"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"golang.org/x/text/message"

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

	workDir, _ := os.Getwd()

	// define handler that serves HTTP requests with the content of the static assets
	staticAssetsDir := filepath.Join(workDir, "/assets")
	fs := http.FileServer(http.Dir(staticAssetsDir))
	// strip the prefix so the path isn't duplicated, which would return an error
	fs = http.StripPrefix("/assets", fs)

	// serve static assets
	r.Get("/assets/*", func(w http.ResponseWriter, r *http.Request) {
		fs.ServeHTTP(w, r)
	})

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		totals, err := dbClient.GetVaccinationTotals(ctx)
		if err != nil {
			fmt.Fprintf(w, "failed to get vaccination totals %v", err)
		}

		ret := IndexPage{
			Pfizer:  totals.Pfizer,
			Moderna: totals.Moderna,
			Janssen: totals.Janssen,
		}

		render(w, "templates/index.html", ret)
	})

	r.Get("/about/", func(w http.ResponseWriter, r *http.Request) {
		render(w, "templates/about.html", "About")
	})

	r.Get("/vaccine/{vaccine}/", func(w http.ResponseWriter, r *http.Request) {
		vaccineSlug := chi.URLParam(r, "vaccine")
		vaccine := store.ManufacturerFromString(vaccineSlug)

		counts, err := dbClient.GetSymptomCounts(ctx, vaccine)
		if err != nil {
			fmt.Fprintf(w, "failed to get symptoms %v", err)
		}

		ret := VaccinePage{
			IsOverview:    true,
			PageTitle:     vaccine.String(),
			TabTitle:      vaccine.String(),
			Vaccine:       vaccine.String(),
			VaccineSlug:   vaccineSlug,
			SymptomCounts: counts,
		}

		render(w, "templates/vaccine.html", ret)
	})

	r.Get("/vaccine/{vaccine}/category/{name}/{sex}/{agemin}/{agemax}/", func(w http.ResponseWriter, r *http.Request) {
		sex := store.SexFromString(chi.URLParam(r, "sex"))

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

		vaccineSlug := chi.URLParam(r, "vaccine")
		vaccine := store.ManufacturerFromString(vaccineSlug)

		categorySlug := chi.URLParam(r, "name")
		categoryName, err := dbClient.GetCategoryName(ctx, categorySlug)
		if err != nil {
			fmt.Fprintf(w, "failed to get category %v", err)
		}

		counts, err := dbClient.GetSymptomCounts(ctx, vaccine)
		if err != nil {
			fmt.Fprintf(w, "failed to get symptoms %v", err)
		}

		results, err := dbClient.GetFilteredResults(ctx, sex, int(ageFloor), int(ageCeil), vaccine, categoryName)
		if err != nil {
			fmt.Fprintf(w, "failed to get results %v", err)
		}

		ret := VaccinePage{
			PageTitle:     vaccine.String(),
			TabTitle:      fmt.Sprintf("%s: %s", vaccine.String(), categoryName),
			Vaccine:       vaccine.String(),
			VaccineSlug:   vaccineSlug,
			SymptomCounts: counts,
			ResultsPage: ResultsPage{
				Vaccine:         vaccine.String(),
				CurrentCategory: categoryName,
				AgeMin:          int(ageFloor),
				AgeMax:          int(ageCeil),
				Sex:             sex.String(),
				Results:         results,
			},
		}

		render(w, "templates/vaccine.html", ret)
	})

	r.Get("/*", func(w http.ResponseWriter, r *http.Request) {
		render(w, "templates/404.html", nil)
	})

	log.Fatal(http.ListenAndServe(":8888", r))

	defer conn.Close(ctx)
}

func funcMap() template.FuncMap {
	p := message.NewPrinter(message.MatchLanguage("en"))

	return template.FuncMap{
		"ellipsis": func(s string) string {
			if len(s) > 100 {
				return s[:100] + "..."
			}
			return s
		},
		"comma": func(strs []string) string {
			return strings.Join(strs, ", ")
		},
		"formatNum": p.Sprint,
	}
}

func render(w http.ResponseWriter, templateName string, ret interface{}) {
	fm := funcMap()
	b, err := ioutil.ReadFile(templateName)
	if err != nil {
		fmt.Fprintf(w, "failed to read template file %v", err)
	}

	t, err := template.New("").Funcs(fm).Parse(string(b))
	if err != nil {
		fmt.Fprintf(w, "failed to parse template %v", err)
	}

	t, err = t.ParseFiles("templates/header.html", "templates/footer.html", "templates/last_updated.html")
	if err != nil {
		fmt.Fprintf(w, "failed to parse partial templates %v", err)
	}

	if err := t.Execute(w, ret); err != nil {
		fmt.Fprintf(w, "failed to execute template %v", err)
	}
}

type VaccinePage struct {
	IsOverview    bool
	PageTitle     string
	TabTitle      string
	Vaccine       string
	VaccineSlug   string
	SymptomCounts []store.SymptomCount
	ResultsPage   ResultsPage
}

type ResultsPage struct {
	Vaccine         string
	CurrentCategory string
	AgeMin          int
	AgeMax          int
	Sex             string
	Results         []store.FilteredResult
}

type IndexPage struct {
	Pfizer  int64
	Moderna int64
	Janssen int64
}
