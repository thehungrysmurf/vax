package store

import "time"

type Store interface {
	InsertReport(report Report) (int64, error)
}

type DB struct {

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

type Symptom struct {
	Name string
	Alias string
}

type Illness string

const(
	Covid19 Illness = "covid19"
)

type Manufacturer string

const(
	Moderna Manufacturer = "moderna"
	Pfizer = "pfizer"
	JohnsonAndJohnson = "janssen"
)

type Vaccine struct {
	Illness Illness
	Manufacturer Manufacturer
}

const InsertReportQuery = `INSERT INTO people (vaers_id, age, sex, notes, reported_at) VALUES ($1, $2, $3, $4, $5);`
const InsertSymptomQuery = ``
const InsertCategoryQuery = ``

func(s *DB) InsertReport(r Report) (int64, error) {
	return 0, nil
}

/*

1. parse reports file, insert into people table
2. parse vaccines file, store in map

2.
option 1
parse symptoms file, for each record
- select people.id corresponding to the source vaers_id - NOW WE HAVE THE PEOPLE ID
- select vaccine id
- populate symptoms lookup and insert if not exists - NOW WE HAVE THE SYMPTOMS ID
- populate people_symptoms NOW WE HAVE BOTH IDS AND WE CAN POPULATE THIS TABLE
- populate symptoms_categories (after lookup in a symptoms<>categories map)
-

map[VaersID]Summary struct {
	SymptomIDs []int
	VaccineID int
}

- parse reports file, insert into people table, add Summary empty to map
- parse vaccines file, insert into vaccines table, lookup Summary from map, set VaccineID
- parse symptoms files, insert into symptoms table, lookup Summary from map, append SymptomID
- range over map, insert into people_symptoms

*/
