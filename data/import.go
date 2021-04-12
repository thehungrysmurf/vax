package data

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

	"github.com/thehungrysmurf/vax/db/store"
)

const Covid19 = "covid19"

type Importer interface {
	Run() error
	ReadReportsFile() error
	ReadVaccinesFile() (map[int64]bool, error)
	ReadSymptomsFile() (map[string]int, error)
}

type CSVImporter struct {
	ReportsFilePath string
	VaccinesFilePath string
	SymptomsFilePath string
	DBClient *store.DB
}

type Summary struct {
	Symptoms []store.Symptom
	VaccineID int
}

func NewCSVImporter(reportsFilePath, vaccinesFilePath, symptomsFilePath string, dbClient *store.DB) CSVImporter {
	return CSVImporter{
		ReportsFilePath:  reportsFilePath,
		VaccinesFilePath: vaccinesFilePath,
		SymptomsFilePath: symptomsFilePath,
		DBClient: dbClient,
	}
}

func (i CSVImporter) Run() error {
	summaryMap := map[int64]*Summary{}
	ctx := context.Background()

	vaccineMap, err := i.ReadVaccinesFile(ctx, summaryMap)
	if err != nil {
		log.Fatalf("failed to read vaccines file: %v", err)
		return err
	}

	if err := i.ReadReportsFile(ctx, summaryMap); err != nil {
		log.Fatalf("failed to read reports file: %v", err)
		return err
	}

	if _, err = i.ReadSymptomsFile(ctx, vaccineMap, summaryMap); err != nil {
		log.Fatalf("failed to read symptoms file: %v", err)
		return err
	}

	// TODO alert if any symptoms with count 50+ are not in category map

	for id, s := range summaryMap {
		log.Printf("complete summary map, vaersID: %v, summary: %+v", id, *s)
	}

	// Populate people_symptoms, symptoms_categories
	for vaersID, summary := range summaryMap {
		for _, symptom := range summary.Symptoms {
			if err := i.DBClient.InsertPeopleSymptom(ctx, vaersID, symptom.ID, summary.VaccineID); err != nil {
				log.Printf("failed to insert people symptoms row for vaers_id: %v, symptom_id: %v, vaccine_id: %v %v", vaersID, symptom.ID, summary.VaccineID, err)
				continue
			}

			for _, cID:= range symptom.CategoryIDs {
				if err := i.DBClient.InsertSymptomCategory(ctx, symptom.ID, cID); err != nil {
					log.Printf("failed to insert symptoms categories row for vaers_id: %v, symptom_id: %v, category_id: %v %v", vaersID, symptom.ID, cID, err)
					continue
				}
			}
		}
	}

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

				var manufacturer store.Manufacturer
				manufacturer = manufacturer.FromString(line[2])
				v := store.Vaccine{
					Illness:      Covid19,
					Manufacturer: manufacturer,
				}

				vaccineID, err := i.DBClient.GetVaccineID(ctx, v)
				if err != nil {
					log.Printf("failed to get vaccine ID for vaers_id %v", line[0])
					continue
				}

				vaersID, err := strconv.ParseInt(line[0], 10, 64)
				if err != nil {
					log.Printf("failed to convert vaers_id %q to int64, skipping row", line[0])
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
				log.Printf("failed to convert vaers_id %q to int64, skipping row", line[0])
				continue
			}

			// Add to DB only if covid19 row
			if _, ok := summaryMap[vaersID]; ok {
				var age float64
				if line[3] != "" {
					age, err = strconv.ParseFloat(line[3], 0)
					if err != nil {
						log.Printf("failed to convert age %q to int64, skipping row", line[3])
						continue
					}
				}

				reportedAt, err := time.Parse("01/02/2006", line[1])
				if err != nil {
					log.Printf("failed to convert reportedAt %q to time format, skipping row", line[1])
					continue
				}

				var s store.Sex
				r := store.Report{
					VaersID:    vaersID,
					Age:        int(age),
					Sex:     s.FromString(line[6]),
					Notes:      line[8],
					ReportedAt: reportedAt,
				}

				if err = i.DBClient.InsertReport(ctx, r); err != nil {
					log.Printf("failed to insert report for vaers_id %s %v", r.VaersID, err)
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
					categories, ok := categoriesMap[s]
					if !ok {
						log.Printf("symptom %s not found in categories map, skipping", s)
						continue
					}

					symptom := store.Symptom{Name:  s}
					if a, ok := aliasesMap[s]; ok {
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
						log.Printf("failed to convert vaers_id %q to int64, skipping row", line[0])
						continue
					}

					if _, ok := summaryMap[vaersID]; !ok {
						log.Printf("failed to fetch summary for vaers_id %s, skipping row", vaersID)
						continue
					}

					var categoryIDs []int
					for _, c := range categories {
						cID, err := i.DBClient.GetCategoryID(ctx, c)
						if err != nil {
							log.Printf("failed to fetch category ID for category %s", c)
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
