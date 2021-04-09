package store

import (
	"github.com/jackc/pgx/v4"
)

type Store interface {
	InsertReport(report Report) (int64, error)
	InsertVaccine(v Vaccine) (int64, error)
	InsertSymptom(s Symptom) (int64, error)
	InsertPeopleSymptoms(vaccineID int, symptomID int64)
	InsertSymptomsCategories(symptomID int, categoryID int64)
	GetCategoryID(category string) (int, error)
}

type DB struct {
	conn *pgx.Conn
}

const InsertReportQuery = `INSERT INTO people (vaers_id, age, sex, notes, reported_at) VALUES ($1, $2, $3, $4, $5);`
const InsertVaccineQuery = `INSERT INTO vaccines (illness, manufacturer) VALUES ($1, $2);`
const InsertSymptomQuery = `INSERT INTO symptoms (symptom, alias) VALUES ($1, $2);`
const InsertPeopleSymptomsQuery = `INSERT INTO people_symptoms(vaers_id, symptom_id, vaccine_id) VALUES ($1, $2, $3);`
const InsertSymptomsCategoryQuery = `INSERT INTO symptoms_categories (symptom_id, category_id) VALUES ($1, $2);`

func NewDB(conn *pgx.Conn) *DB {
	return &DB{
		conn: conn,
	}
}
func(d *DB) InsertReport(r Report) (int64, error) {
	return 1, nil
}

func(d *DB) InsertVaccine(v Vaccine) (int64, error) {
	return 2, nil
}

func(d *DB) InsertSymptom(s Symptom) (int64, error) {
	return 3, nil
}

func(d *DB) InsertPeopleSymptoms(vaccineID int, symptomID int64) error {
	return nil
}

func(d *DB) InsertSymptomsCategories(symptomID int64, categoryID int) error {
	return nil
}

func(d *DB) GetCategoryID(category string) (int, error) {
	return 12, nil
}
