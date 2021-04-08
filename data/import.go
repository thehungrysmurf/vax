package data

import (
	"bufio"
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
	SymptomIDs []int64
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
	summaryMap := map[int64]Summary{}

	err := i.ReadReportsFile(summaryMap)
	if err != nil {
		log.Fatalf("failed to read reports file: %v", err)
		return err
	}

	vaccineMap, err := i.ReadVaccinesFile(summaryMap)
	if err != nil {
		log.Fatalf("failed to read vaccines file: %v", err)
		return err
	}

	fmt.Printf("size of vaccine map: %v \n", len(vaccineMap))

	_, err = i.ReadSymptomsFile(vaccineMap, summaryMap)
	if err != nil {
		log.Fatalf("failed to read symptoms file: %v", err)
		return err
	}

	// range over map, insert into people_symptoms
	for vaers_id, summary := range summaryMap {
		for _, symptomID := range summary.SymptomIDs {
			if err := i.DBClient.InsertPeopleSymptoms(summary.VaccineID, symptomID); err != nil {
				log.Printf("failed to insert people symptoms row for vaers_id %v", vaers_id)
				continue
			}
		}
	}

	return nil
}

// parse reports file, insert into people table, add Summary empty to map
func (i CSVImporter) ReadReportsFile(summaryMap map[int64]Summary) error {
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

			var g store.Gender
			r := store.Report{
				VaersID:    vaersID,
				Age:        int(age),
				Gender:     g.FromString(line[6]),
				Notes:      line[8],
				ReportedAt: reportedAt,
			}

			if _, err := i.DBClient.InsertReport(r); err != nil {
				log.Printf("failed to insert report for vaers_id %s", r.VaersID)
				continue
			}

			summaryMap[r.VaersID] = Summary{}
		}
	}

	log.Printf("finished reading reports file, read %d lines", linesRead)
	return nil
}

// parse vaccines file, insert into vaccines table, lookup Summary from map, set VaccineID
func (i CSVImporter) ReadVaccinesFile(summaryMap map[int64]Summary) (map[int64]bool, error) {
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
			if line[1] == "COVID19" {
				id, err := strconv.ParseInt(line[0], 10, 64)
				if err != nil {
					log.Printf("failed to convert ID %s to int: %v", line[0], err)
					continue
				}
				vaccineMap[id] = true

				var m store.Manufacturer
				v := store.Vaccine{
					Illness:      store.Covid19,
					Manufacturer: m.FromString(line[2]),
				}
				vaccineID, err := i.DBClient.InsertVaccine(v)
				if err != nil {
					log.Printf("failed to insert vaccine for vaers_id %s, skipping row", line[0])
					continue
				}

				vaersID, err := strconv.ParseInt(line[0], 10, 64)
				if err != nil {
					log.Printf("failed to convert vaers_id %q to int64, skipping row", line[0])
					continue
				}

				s, ok := summaryMap[vaersID]
				if !ok {
					log.Printf("failed to fetch summary for vaers_id %s, skipping row", vaersID)
					continue
				}

				s.VaccineID = int(vaccineID)
			}
		}
	}

	log.Printf("finished reading vaccines file, read %d lines", linesRead)
	return vaccineMap, nil
}

// parse symptoms files, insert into symptoms table, lookup Summary from map, append SymptomID
func (i *CSVImporter) ReadSymptomsFile(vaccineMap map[int64]bool, summaryMap map[int64]Summary) (map[string]int, error) {
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

			for _, s:= range symptoms {
				s := store.Symptom{
					Name:  s,
					Alias: "TBD",
				}
				symptomID, err := i.DBClient.InsertSymptom(s)
				if err != nil {
					log.Printf("failed to insert symptom for vaers_id %s, skipping row", line[0])
					continue
				}

				vaersID, err := strconv.ParseInt(line[0], 10, 64)
				if err != nil {
					log.Printf("failed to convert vaers_id %q to int64, skipping row", line[0])
					continue
				}

				summary, ok := summaryMap[vaersID]
				if !ok {
					log.Printf("failed to fetch summary for vaers_id %s, skipping row", vaersID)
					continue
				}

				summary.SymptomIDs = append(summary.SymptomIDs, symptomID)
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
