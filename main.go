package main

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
)

func main() {
	symptomsFilePath := "/Users/thesmurf/Documents/2021VAERSData/2021VAERSSYMPTOMS.csv"
	vaxFilePath := "/Users/thesmurf/Documents/2021VAERSData/2021VAERSVAX.csv"

	vaccineMap, err := readVaccineFile(vaxFilePath)
	if err != nil {
		log.Fatalf("failed to read vaccine file: %v", err)
	}

	fmt.Printf("size of vaccine map: %v \n", len(vaccineMap))

	symptomsMap, err := readSymptomsFile(symptomsFilePath, vaccineMap)
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

	for _, s := range sorted {
		if symptomsMap[s] > 10 {
			fmt.Printf("%v %v\n", s, symptomsMap[s])
		}
	}
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

	// line: [VAERS_ID SYMPTOM1 SYMPTOMVERSION1 SYMPTOM2 SYMPTOMVERSION2 SYMPTOM3 SYMPTOMVERSION3 SYMPTOM4 SYMPTOMVERSION4 SYMPTOM5 SYMPTOMVERSION5]
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

func readVaccineFile(vaccineFilePath string) (map[int64]bool, error) {
	csvFile, err := os.Open(vaccineFilePath)
	if err != nil {
		log.Printf("failed to parse csv file: %v", err)
		return nil, err
	}

	reader := csv.NewReader(bufio.NewReader(csvFile))
	linesRead := 0
	vaccineMap := map[int64]bool{}

	// [1095070 COVID19 PFIZER\BIONTECH EL9265 1 IM LA COVID19 (COVID19 (PFIZER-BIONTECH))]
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
