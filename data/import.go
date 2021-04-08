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
}

func NewCSVImporter(reportsFilePath, vaccinesFilePath, symptomsFilePath string) CSVImporter {
	return CSVImporter{
		ReportsFilePath:  reportsFilePath,
		VaccinesFilePath: vaccinesFilePath,
		SymptomsFilePath: symptomsFilePath,
	}
}

func (i *CSVImporter) Run() error {
	err := i.ReadReportsFile()
	if err != nil {
		log.Fatalf("failed to read reports file: %v", err)
		return err
	}

	vaccineMap, err := i.ReadVaccinesFile()
	if err != nil {
		log.Fatalf("failed to read vaccines file: %v", err)
		return err
	}

	fmt.Printf("size of vaccine map: %v \n", len(vaccineMap))

	symptomsMap, err := i.ReadSymptomsFile(vaccineMap)
	if err != nil {
		log.Fatalf("failed to read symptoms file: %v", err)
		return err
	}

	sorted := make([]string, 0, len(symptomsMap))
	for s := range symptomsMap {
		sorted = append(sorted, s)
	}
	sort.Slice(sorted, func(i, j int) bool {
		return symptomsMap[sorted[i]] > symptomsMap[sorted[j]]
	})

	// for _, s := range sorted {
	// 	if symptomsMap[s] > 10 {
	// 		fmt.Printf("%v %v\n", s, symptomsMap[s])
	// 	}
	// }

	return nil
}

func (i *CSVImporter) ReadReportsFile() error {
	csvFile, err := os.Open(i.ReportsFilePath)
	if err != nil {
		log.Printf("failed to parse csv file: %v", err)
		return err
	}

	reader := csv.NewReader(bufio.NewReader(csvFile))
	linesRead := 0

	for {
		line, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			log.Printf("failed to read from csv file: %v", err)
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

			_ = store.Report{
				VaersID:    vaersID,
				Age:        int(age),
				Gender:     store.Gender(line[6]),
				Notes:      line[8],
				ReportedAt: reportedAt,
			}
		}
	}

	fmt.Printf("finished reading reports file, read %d lines \n", linesRead)
	return nil
}

func (i *CSVImporter) ReadVaccinesFile() (map[int64]bool, error) {
	csvFile, err := os.Open(i.VaccinesFilePath)
	if err != nil {
		log.Printf("failed to parse csv file: %v", err)
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
			log.Printf("failed to read from csv file: %v", err)
			return nil, err
		}
		linesRead++

		if linesRead > 1 {
			if line[1] == "COVID19" {
				id, err := strconv.ParseInt(line[0], 10, 64)
				if err != nil {
					return nil, err
				}
				vaccineMap[id] = true
			}
		}
	}

	fmt.Printf("finished reading vaccine file, read %d lines \n", linesRead)
	return vaccineMap, nil
}

func (i *CSVImporter) ReadSymptomsFile(vaccineMap map[int64]bool) (map[string]int, error) {
	csvFile, err := os.Open(i.SymptomsFilePath)
	if err != nil {
		log.Printf("failed to parse csv file: %v", err)
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
			log.Printf("failed to read from csv file: %v", err)
			return nil, err
		}
		linesRead++

		if linesRead > 1 {
			id, err := strconv.ParseInt(line[0], 10, 64)
			if err != nil {
				return nil, err
			}
			if _, ok := vaccineMap[id]; ok {
				symptomsToAdd := []string{line[1], line[3], line[5], line[7], line[9]}
				loadSymptoms(symptomsMap, symptomsToAdd)
			}
		}
	}
	fmt.Printf("finished reading symptoms file, read %d lines \n", linesRead)
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
