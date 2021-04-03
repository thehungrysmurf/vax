package main

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

	"github.com/thehungrysmurf/vax/config"

	"github.com/jackc/pgx/v4"
	"github.com/joeshaw/envdecode"
)

func main() {
	var cfg config.Config
	err := envdecode.Decode(&cfg)
	if err != nil {
		log.Fatalf("failed to read config: %v", err)
	}

	err = readReportsFile(cfg.ReportsFilePath)
	if err != nil {
		log.Fatalf("failed to read reports file: %v", err)
	}

	vaccineMap, err := readVaccinesFile(cfg.VaccinesFilePath)
	if err != nil {
		log.Fatalf("failed to read vaccines file: %v", err)
	}

	fmt.Printf("size of vaccine map: %v \n", len(vaccineMap))

	symptomsMap, err := readSymptomsFile(cfg.SymptomsFilePath, vaccineMap)
	if err != nil {
		log.Fatalf("failed to read symptoms file: %v", err)
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

	ctx := context.Background()
	conn, err := pgx.Connect(ctx, cfg.DatabaseURI)
	if err != nil {
		log.Fatal("failed to connect to database")
	}

	defer conn.Close(ctx)
}

func readSymptomsFile(symptomsFilePath string, vaccineMap map[int64]bool) (map[string]int, error) {
	csvFile, err := os.Open(symptomsFilePath)
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

func readVaccinesFile(vaccinesFilePath string) (map[int64]bool, error) {
	csvFile, err := os.Open(vaccinesFilePath)
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

type Gender string

const(
	Male Gender = "M"
	Female = "F"
	Unknown = "U"
)

type Report struct {
	VaersID int64
	Age int
	Gender Gender
	Notes string
	ReportedAt time.Time
}

func readReportsFile(reportsFilePath string) (error) {
	csvFile, err := os.Open(reportsFilePath)
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

			_ = Report{
				VaersID:    vaersID,
				Age:        int(age),
				Gender:     Gender(line[6]),
				Notes:      line[8],
				ReportedAt: reportedAt,
			}
		}
	}

	fmt.Printf("finished reading reports file, read %d lines \n", linesRead)
	return nil
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
