package store

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v4"
	"github.com/thehungrysmurf/vax/data"
)

type Store interface {
	InsertVaccinationTotals(ctx context.Context, totals VaccinationTotals) error
	InsertReport(ctx context.Context, r Report) (int64, error)
	InsertSymptom(ctx context.Context, s Symptom) (int64, error)
	InsertPeopleSymptoms(ctx context.Context, vaxID int, symID int64)
	InsertSymptomsCategories(ctx context.Context, symID int, catID int64)
	GetVaccinationTotals(ctx context.Context) (VaccinationTotals, error)
	GetVaccineID(ctx context.Context, v Vaccine) (int, error)
	GetCategoryID(ctx context.Context, cat string) (int, error)
	GetCategoryName(ctx context.Context, catSlug string) (string, error)
	GetSymptomCounts(ctx context.Context, manufacturer Manufacturer) ([]CategoryCount, error)
	GetFilteredResults(ctx context.Context, sex Sex, ageFloor, ageCeiling int, manufacturer Manufacturer, categoryName string) ([]FilteredResult, error)
}

type DB struct {
	conn *pgx.Conn
}

type VaccinationTotals struct {
	Pfizer  int64
	Moderna int64
	Janssen int64
}

func NewDB(conn *pgx.Conn) *DB {
	return &DB{
		conn: conn,
	}
}

const InsertVaccinationTotalsQuery = `INSERT INTO vaccination_totals (pfizer, moderna, janssen, updated_at) values ($1, $2, $3, $4)`

func (d *DB) InsertVaccinationTotals(ctx context.Context, totals VaccinationTotals) error {
	_, err := d.conn.Exec(ctx, InsertVaccinationTotalsQuery, totals.Pfizer, totals.Moderna, totals.Janssen, time.Now())
	return err
}

const SelectVaccinationTotalsQuery = `SELECT pfizer, moderna, janssen FROM vaccination_totals ORDER BY updated_at DESC LIMIT 1`

func (d *DB) GetVaccinationTotals(ctx context.Context) (VaccinationTotals, error) {
	var vt VaccinationTotals
	err := d.conn.QueryRow(ctx, SelectVaccinationTotalsQuery).Scan(&vt.Pfizer, &vt.Moderna, &vt.Janssen)
	return vt, err
}

const InsertReportQuery = `INSERT INTO people (vaers_id, age, sex, notes, reported_at) VALUES ($1, $2, $3, $4, $5);`

func (d *DB) InsertReport(ctx context.Context, r Report) error {
	_, err := d.conn.Exec(ctx, InsertReportQuery, r.VaersID, r.Age, r.Sex, r.Notes, r.ReportedAt)
	return err
}

const SelectVaccineQuery = `SELECT id FROM vaccines WHERE illness = $1 AND manufacturer = $2;`

func (d *DB) GetVaccineID(ctx context.Context, v Vaccine) (int, error) {
	var id int
	err := d.conn.QueryRow(ctx, SelectVaccineQuery, v.Illness, v.Manufacturer).Scan(&id)
	return id, err
}

const SelectSymptomQuery = `SELECT id FROM symptoms WHERE name = $1`
const InsertSymptomQuery = `INSERT INTO symptoms (name, alias) VALUES ($1, $2) RETURNING id;`

func (d *DB) InsertSymptom(ctx context.Context, s Symptom) (int64, error) {
	var id int64
	err := d.conn.QueryRow(ctx, SelectSymptomQuery, s.Name).Scan(&id)
	if err == pgx.ErrNoRows {
		err = d.conn.QueryRow(ctx, InsertSymptomQuery, s.Name, s.Alias).Scan(&id)
	}

	return id, err
}

const InsertPeopleSymptomQuery = `INSERT INTO people_symptoms(vaers_id, symptom_id, vaccine_id) VALUES ($1, $2, $3) ON CONFLICT DO NOTHING;`

func (d *DB) InsertPeopleSymptom(ctx context.Context, vaersID, symID int64, vaxID int) error {
	_, err := d.conn.Exec(ctx, InsertPeopleSymptomQuery, vaersID, symID, vaxID)
	return err
}

const InsertSymptomCategoryQuery = `INSERT INTO symptoms_categories (symptom_id, category_id) VALUES ($1, $2) ON CONFLICT DO NOTHING;`

func (d *DB) InsertSymptomCategory(ctx context.Context, symID int64, catID int) error {
	_, err := d.conn.Exec(ctx, InsertSymptomCategoryQuery, symID, catID)
	return err
}

const SelectCategoryIDQuery = `SELECT id FROM categories WHERE name = $1`

func (d *DB) GetCategoryID(ctx context.Context, cat string) (int, error) {
	var id int
	err := d.conn.QueryRow(ctx, SelectCategoryIDQuery, cat).Scan(&id)
	return id, err
}

const SelectCategoryNameQuery = `SELECT name FROM categories WHERE slug = $1`

func (d *DB) GetCategoryName(ctx context.Context, catSlug string) (string, error) {
	var name string
	err := d.conn.QueryRow(ctx, SelectCategoryNameQuery, catSlug).Scan(&name)
	return name, err
}

type CategoryCount struct {
	Category     string `db:"category"`
	CategorySlug string `db:"slug"`
	Count        int64  `db:"count"`
}

const SelectCategoryCountsQuery = `SELECT c.name as category, c.slug as slug, count(ps.vaers_id) as count FROM categories c
JOIN symptoms_categories sc ON c.id = sc.category_id
JOIN symptoms s ON s.id = sc.symptom_id
JOIN people_symptoms ps ON ps.symptom_id = s.id
JOIN vaccines v ON v.id = ps.vaccine_id
WHERE v.manufacturer = $1
AND c.slug != 'errors-by-medical-staff'
GROUP BY c.name, c.slug;`

func (d *DB) GetCategoryCounts(ctx context.Context, manufacturer Manufacturer) ([]CategoryCount, error) {
	var counts []CategoryCount
	rows, err := d.conn.Query(ctx, SelectCategoryCountsQuery, manufacturer)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		sc := CategoryCount{}
		if err := rows.Scan(&sc.Category, &sc.CategorySlug, &sc.Count); err != nil {
			return nil, fmt.Errorf("failed to scan result: %v", err)
		}
		counts = append(counts, sc)
	}

	log.Printf("--> Found category counts: %#+v", counts)
	return counts, nil
}

type FilteredResult struct {
	Age        int      `db:"age"`
	ReportedAt string   `db:"reported_at"`
	Notes      string   `db:"notes"`
	Symptoms   []string `db:"symptoms"`
}

const SelectFilteredResultsQuery = `SELECT p.age as age, p.reported_at as reported_at, p.notes as notes, json_agg(s.name) as symptoms FROM people p
JOIN people_symptoms ps ON p.vaers_id = ps.vaers_id
JOIN symptoms s ON s.id = ps.symptom_id
JOIN symptoms_categories sc ON sc.symptom_id = s.id
JOIN categories c ON c.id = sc.category_id
JOIN vaccines v ON v.id = ps.vaccine_id
WHERE p.sex = $1 
AND p.age BETWEEN $2 AND $3 
AND v.manufacturer = $4
AND c.name = $5
GROUP BY p.age, p.notes, p.reported_at
ORDER BY p.age, p.reported_at;
`

func (d *DB) GetFilteredResults(ctx context.Context, sex Sex, ageMin, ageMax int, manufacturer Manufacturer, category string) ([]FilteredResult, error) {
	var results []FilteredResult
	rows, err := d.conn.Query(ctx, SelectFilteredResultsQuery, sex, ageMin, ageMax, manufacturer, category)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		fr := FilteredResult{}
		var reportedAt time.Time
		if err := rows.Scan(&fr.Age, &reportedAt, &fr.Notes, &fr.Symptoms); err != nil {
			return nil, fmt.Errorf("failed to scan result: %v", err)
		}
		fr.ReportedAt = reportedAt.Format("2006-01-02")

		// Replace symptoms with their plain English synonyms, if they exist
		for i, sym := range fr.Symptoms {
			if alias, ok := data.AliasesMap[sym]; ok {
				fr.Symptoms[i] = alias
			}
		}

		results = append(results, fr)
	}

	log.Printf("--> Found %d filtered results.", len(results))
	return results, nil
}

type SymptomCount struct {
	Symptom  string
	Count    int64
	Category string
}

const SelectSymptomCountQuery = `
SELECT s.name AS symptom, c.name AS category, count(ps.vaers_id) AS count FROM categories c
JOIN symptoms_categories sc ON c.id = sc.category_id
JOIN symptoms s ON s.id = sc.symptom_id
JOIN people_symptoms ps ON ps.symptom_id = s.id
JOIN vaccines v ON v.id = ps.vaccine_id
WHERE v.manufacturer = $1 AND c.slug != 'errors-by-medical-staff'
GROUP BY s.name, c.name ORDER BY count(ps.vaers_id) DESC
LIMIT 30;
`

func (d *DB) GetSymptomCounts(ctx context.Context, manufacturer Manufacturer) ([]SymptomCount, error) {
	var results []SymptomCount
	rows, err := d.conn.Query(ctx, SelectSymptomCountQuery, manufacturer)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		sc := SymptomCount{}
		if err := rows.Scan(&sc.Symptom, &sc.Category, &sc.Count); err != nil {
			return nil, fmt.Errorf("failed to scan result: %v", err)
		}

		// Replace symptom with its plain English synonyms, if it exists
		if alias, ok := data.AliasesMap[sc.Symptom]; ok {
			sc.Symptom = alias
		}

		results = append(results, sc)
	}

	log.Printf("--> Found symptom counts: %#+v", results)
	return results, nil
}
