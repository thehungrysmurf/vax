package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"sort"
	"strings"
)

func main() {
	csvFile, err := os.Open("/Users/thesmurf/Documents/2021VAERSData/2021VAERSSYMPTOMS.csv")
	if err != nil {
		log.Printf("failed to parse csv file: %v", err)
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
			log.Fatalf("failed to read from csv file: %v", err)
		}
		linesRead++

		if linesRead > 1 {
			symptomsToAdd := []string{line[1], line[3], line[5], line[7], line[9]}
			loadSymptoms(symptomsMap, symptomsToAdd)
		}
	}
	fmt.Printf("all done, read %d lines \n", linesRead)

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
