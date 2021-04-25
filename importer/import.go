package importer

import (
	"bufio"
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/thehungrysmurf/vax/data"
	"github.com/thehungrysmurf/vax/db/store"
)

const Covid19 = "covid19"

type Importer interface {
	Run() error
	ReadVaccinationTotalsFile() error
	ReadReportsFile() error
	ReadVaccinesFile() (map[int64]bool, error)
	ReadSymptomsFile() (map[string]int, error)
}

type CSVImporter struct {
	VaccinationTotalsFilePath string
	ReportsFilePath           string
	VaccinesFilePath          string
	SymptomsFilePath          string
	DBClient                  *store.DB
}

type Summary struct {
	Symptoms  []store.Symptom
	VaccineID int
}

func NewCSVImporter(vaccinationTotalsFilePath, reportsFilePath, vaccinesFilePath, symptomsFilePath string, dbClient *store.DB) CSVImporter {
	return CSVImporter{
		VaccinationTotalsFilePath: vaccinationTotalsFilePath,
		ReportsFilePath:           reportsFilePath,
		VaccinesFilePath:          vaccinesFilePath,
		SymptomsFilePath:          symptomsFilePath,
		DBClient:                  dbClient,
	}
}

func (i CSVImporter) Run() error {
	summaryMap := map[int64]*Summary{}
	ctx := context.Background()

	err := i.ReadVaccinationTotalsFile(ctx)
	if err != nil {
		log.Fatalf("failed to read vaccination totals file: %v", err)
		return err
	}

	vaccineMap, err := i.ReadVaccinesFile(ctx, summaryMap)
	if err != nil {
		log.Fatalf("failed to read vaccines file: %v", err)
		return err
	}

	if err := i.ReadReportsFile(ctx, summaryMap); err != nil {
		log.Fatalf("failed to read reports file: %v", err)
		return err
	}

	symptomsMap, err := i.ReadSymptomsFile(ctx, vaccineMap, summaryMap)
	if err != nil {
		log.Fatalf("failed to read symptoms file: %v", err)
		return err
	}

	for s, count := range symptomsMap {
		if count >= 50 {
			if _, ok := data.CategoriesMap[s]; !ok {
				log.Printf("!! symptom %s has been reported %v times and needs to be categorized !!", s, count)
			}
		}
	}

	// for id, s := range summaryMap {
	// 	log.Printf("complete summary map, vaersID: %v, summary: %+v", id, *s)
	// }

	// Populate people_symptoms, symptoms_categories
	for vaersID, summary := range summaryMap {
		for _, symptom := range summary.Symptoms {
			if err := i.DBClient.InsertPeopleSymptom(ctx, vaersID, symptom.ID, summary.VaccineID); err != nil {
				log.Printf("failed to insert people symptoms row for vaers_id: %v, symptom_id: %v, vaccine_id: %v %v", vaersID, symptom.ID, summary.VaccineID, err)
				continue
			}

			for _, cID := range symptom.CategoryIDs {
				if err := i.DBClient.InsertSymptomCategory(ctx, symptom.ID, cID); err != nil {
					log.Printf("failed to insert symptoms categories row for vaers_id: %v, symptom_id: %v, category_id: %v %v", vaersID, symptom.ID, cID, err)
					continue
				}
			}
		}
	}

	return nil
}

// Parse vaccination totals file, insert into vaccination_totals table
func (i CSVImporter) ReadVaccinationTotalsFile(ctx context.Context) error {
	csvFile, err := os.Open(i.VaccinationTotalsFilePath)
	if err != nil {
		log.Printf("failed to open csv vaccination totals file: %v", err)
		return err
	}

	reader := csv.NewReader(bufio.NewReader(csvFile))
	linesRead := 0
	var vaxTotal store.VaccinationTotals

	for {
		line, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			log.Printf("failed to read from vaccination totals csv file: %v", err)
			return err
		}
		linesRead++

		if linesRead > 1 {
			if line[0] == "United States" {
				if line[2] == `Pfizer/BioNTech` {
					pfizerTotal, err := strconv.ParseInt(line[3], 10, 64)
					if err != nil {
						log.Printf("failed to convert pfizer count %s to int: %v", line[0], err)
						continue
					}
					vaxTotal.Pfizer = pfizerTotal
				}

				if line[2] == "Moderna" {
					modernaTotal, err := strconv.ParseInt(line[3], 10, 64)
					if err != nil {
						log.Printf("failed to convert moderna count %s to int: %v", line[0], err)
						continue
					}
					vaxTotal.Moderna = modernaTotal
				}

				if line[2] == "Johnson&Johnson" {
					janssenTotal, err := strconv.ParseInt(line[3], 10, 64)
					if err != nil {
						log.Printf("failed to convert j&j count %s to int: %v", line[0], err)
						continue
					}
					vaxTotal.Janssen = janssenTotal
				}
			}
		}
	}

	if err := i.DBClient.InsertVaccinationTotals(ctx, vaxTotal); err != nil {
		log.Printf("failed to insert latest vaccination totals: %v", err)
	}

	log.Printf("finished reading vaccination totals file, read %d lines", linesRead)
	return nil
}

// Parse vaccines file, insert into vaccines table, set VaccineID in summary map
func (i CSVImporter) ReadVaccinesFile(ctx context.Context, summaryMap map[int64]*Summary) (map[int64]bool, error) {
	csvFile, err := os.Open(i.VaccinesFilePath)
	if err != nil {
		log.Printf("failed to open csv vaccines file: %v", err)
		return nil, err
	}

	reader := csv.NewReader(bufio.NewReader(csvFile))
	linesRead := 0
	vaccineMap := map[int64]bool{}

	for {
		line, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			log.Printf("failed to read from vaccines csv file: %v", err)
			return nil, err
		}
		linesRead++

		if linesRead > 1 {
			if strings.ToLower(line[1]) == Covid19 {
				id, err := strconv.ParseInt(line[0], 10, 64)
				if err != nil {
					log.Printf("failed to convert ID %s to int: %v", line[0], err)
					continue
				}
				vaccineMap[id] = true

				manufacturer := store.ManufacturerFromString(line[2])

				// Ignore unknown vaccines
				if manufacturer != store.Pfizer && manufacturer != store.Moderna && manufacturer != store.Janssen {
					continue
				}

				v := store.Vaccine{
					Illness:      Covid19,
					Manufacturer: manufacturer,
				}

				vaccineID, err := i.DBClient.GetVaccineID(ctx, v)
				if err != nil {
					log.Printf("failed to get vaccine ID for vaers_id %v: %v", line[0], err)
					continue
				}

				vaersID, err := strconv.ParseInt(line[0], 10, 64)
				if err != nil {
					log.Printf("failed to convert vaers_id %q to int64, skipping row: %v", line[0], err)
					continue
				}

				summaryMap[vaersID] = &Summary{VaccineID: vaccineID}
			}
		}
	}

	log.Printf("finished reading vaccines file, read %d lines", linesRead)
	return vaccineMap, nil
}

// Parse reports file, insert into people table
func (i CSVImporter) ReadReportsFile(ctx context.Context, summaryMap map[int64]*Summary) error {
	csvFile, err := os.Open(i.ReportsFilePath)
	if err != nil {
		log.Printf("failed to open reports csv file: %v", err)
		return err
	}

	reader := csv.NewReader(bufio.NewReader(csvFile))
	linesRead := 0

	for {
		line, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			log.Printf("failed to read from reports csv file: %v", err)
			return err
		}
		linesRead++

		if linesRead > 1 {
			vaersID, err := strconv.ParseInt(line[0], 10, 64)
			if err != nil {
				log.Printf("failed to convert vaers_id %q to int64, skipping row: %v", line[0], err)
				continue
			}

			// Add to DB only if covid19 row
			if _, ok := summaryMap[vaersID]; ok {
				var age float64
				if line[3] != "" {
					age, err = strconv.ParseFloat(line[3], 0)
					if err != nil {
						log.Printf("failed to convert age %q to int64, skipping row: %v", line[3], err)
						continue
					}
				}

				reportedAt, err := time.Parse("01/02/2006", line[1])
				if err != nil {
					log.Printf("failed to convert reportedAt %q to time format, skipping row: %v", line[1], err)
					continue
				}

				r := store.Report{
					VaersID:    vaersID,
					Age:        int(age),
					Sex:        store.SexFromString(line[6]),
					Notes:      line[8],
					ReportedAt: reportedAt,
				}

				if err = i.DBClient.InsertReport(ctx, r); err != nil {
					log.Printf("failed to insert report for vaers_id %v: %v", r.VaersID, err)
					continue
				}
			}
		}
	}

	log.Printf("finished reading reports file, read %d lines", linesRead)
	return nil
}

// Parse symptoms file, insert into symptoms table, lookup categories for symptom and populate Symptoms in summary map
func (i *CSVImporter) ReadSymptomsFile(ctx context.Context, vaccineMap map[int64]bool, summaryMap map[int64]*Summary) (map[string]int, error) {
	csvFile, err := os.Open(i.SymptomsFilePath)
	if err != nil {
		log.Printf("failed to open symptoms csv file: %v", err)
		return nil, err
	}

	reader := csv.NewReader(bufio.NewReader(csvFile))
	linesRead := 0
	symptomsMap := map[string]int{}

	for {
		line, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			log.Printf("failed to read from symptoms csv file: %v", err)
			return nil, err
		}
		linesRead++

		if linesRead > 1 {
			id, err := strconv.ParseInt(line[0], 10, 64)
			if err != nil {
				return nil, err
			}

			symptoms := []string{line[1], line[3], line[5], line[7], line[9]}

			if _, ok := vaccineMap[id]; ok {
				symptomsToAdd := symptoms
				loadSymptoms(symptomsMap, symptomsToAdd)
			}

			for _, s := range symptoms {
				if s != "" {
					s = strings.ToLower(s)
					categories, ok := data.CategoriesMap[s]
					if !ok {
						log.Printf("symptom %s not found in categories map, skipping", s)
						continue
					}

					symptom := store.Symptom{Name: s}
					if a, ok := data.AliasesMap[s]; ok {
						symptom.Alias = a
					}

					sID, err := i.DBClient.InsertSymptom(ctx, symptom)
					if err != nil {
						log.Printf("failed to insert symptom %s for vaers_id %s, skipping row %v", s, line[0], err)
						continue
					}
					symptom.ID = sID

					vaersID, err := strconv.ParseInt(line[0], 10, 64)
					if err != nil {
						log.Printf("failed to convert vaers_id %q to int64, skipping row: %v", line[0], err)
						continue
					}

					if _, ok := summaryMap[vaersID]; !ok {
						log.Printf("failed to fetch summary for vaers_id %v, skipping row", vaersID)
						continue
					}

					var categoryIDs []int
					for _, c := range categories {
						cID, err := i.DBClient.GetCategoryID(ctx, c)
						if err != nil {
							log.Printf("failed to fetch category ID for category %s: %v", c, err)
						}
						categoryIDs = append(categoryIDs, cID)
					}

					symptom.CategoryIDs = categoryIDs
					summaryMap[vaersID].Symptoms = append(summaryMap[vaersID].Symptoms, symptom)
				}
			}
		}
	}

	log.Printf("finished reading symptoms file, read %d lines \n", linesRead)
	return symptomsMap, nil
}

func loadSymptoms(symptomsMap map[string]int, symptomsToAdd []string) map[string]int {
	for _, s := range symptomsToAdd {
		if s != "" {
			s = strings.ToLower(s)
			if _, ok := symptomsMap[s]; ok {
				symptomsMap[s]++
			} else {
				symptomsMap[s] = 1
			}
		}
	}
	return symptomsMap
}

func sortBySymptomCount(symptomsMap map[string]int) {
	sorted := make([]string, 0, len(symptomsMap))
	for s := range symptomsMap {
		sorted = append(sorted, s)
	}
	sort.Slice(sorted, func(i, j int) bool {
		return symptomsMap[sorted[i]] > symptomsMap[sorted[j]]
	})
	for _, s := range sorted {
		if symptomsMap[s] > 10 {
			fmt.Printf("%v %v\n", s, symptomsMap[s])
		}
	}
}
