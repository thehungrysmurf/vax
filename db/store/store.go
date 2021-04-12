package store

import (
	"context"

	"github.com/jackc/pgx/v4"
)

type Store interface {
	InsertReport(r Report) (int64, error)
	InsertSymptom(s Symptom) (int64, error)
	InsertPeopleSymptoms(vaxID int, symID int64)
	InsertSymptomsCategories(symID int, catID int64)
	GetVaccineID(ctx context.Context, v Vaccine) (int, error)
	GetCategoryID(cat string) (int, error)
}

type DB struct {
	conn *pgx.Conn
}

func NewDB(conn *pgx.Conn) *DB {
	return &DB{
		conn: conn,
	}
}

const InsertReportQuery = `INSERT INTO people (vaers_id, age, sex, notes, reported_at) VALUES ($1, $2, $3, $4, $5);`

func(d *DB) InsertReport(ctx context.Context, r Report) error {
	_, err := d.conn.Exec(ctx, InsertReportQuery, r.VaersID, r.Age, r.Sex, r.Notes, r.ReportedAt)
	return err
}

const SelectVaccineQuery = `SELECT id FROM vaccines WHERE illness = $1 AND manufacturer = $2;`

func(d *DB) GetVaccineID(ctx context.Context, v Vaccine) (int, error) {
	var id int
	err := d.conn.QueryRow(ctx, SelectVaccineQuery, v.Illness, v.Manufacturer).Scan(&id)
	return id, err
}

const SelectSymptomQuery = `SELECT id FROM symptoms WHERE name = $1`
const InsertSymptomQuery = `INSERT INTO symptoms (name, alias) VALUES ($1, $2) RETURNING id;`

func(d *DB) InsertSymptom(ctx context.Context, s Symptom) (int64, error) {
	var id int64
	err := d.conn.QueryRow(ctx, SelectSymptomQuery, s.Name).Scan(&id)
	if err == pgx.ErrNoRows {
		err = d.conn.QueryRow(ctx, InsertSymptomQuery, s.Name, s.Alias).Scan(&id)
	}

	return id, err
}

const InsertPeopleSymptomQuery = `INSERT INTO people_symptoms(vaers_id, symptom_id, vaccine_id) VALUES ($1, $2, $3);`

func(d *DB) InsertPeopleSymptom(ctx context.Context, vaersID, symID int64, vaxID int) error {
	_, err := d.conn.Exec(ctx, InsertPeopleSymptomQuery, vaersID, symID, vaxID)
	return err
}

const InsertSymptomCategoryQuery = `INSERT INTO symptoms_categories (symptom_id, category_id) VALUES ($1, $2);`

func(d *DB) InsertSymptomCategory(ctx context.Context, symID int64, catID int) error {
	_, err := d.conn.Exec(ctx, InsertSymptomCategoryQuery, symID, catID)
	return err
}

const SelectCategoryQuery = `SELECT id FROM categories WHERE name = $1`

func(d *DB) GetCategoryID(ctx context.Context, cat string) (int, error) {
	var id int
	err := d.conn.QueryRow(ctx, SelectCategoryQuery, cat).Scan(&id)
	return id, err
}

const SelectReportsPerVaccine = `SELECT c.name, count(ps.vaers_id) FROM categories c
JOIN symptoms_categories sc ON c.id = sc.category_id
JOIN symptoms s ON s.id = sc.symptom_id
JOIN people_symptoms ps ON ps.symptom_id = s.id
JOIN vaccines v ON v.id = ps.vaccine_id
WHERE v.manufacturer = $1
GROUP BY c.name;`

func(d *DB) GetVaccineReports(ctx context.Context, manufacturer Manufacturer) error {
	return nil
}

/* vaccine page
For vaccine X, get me the count of all the people who reported symptoms under each category in the categories table:

select c.name, count(ps.vaers_id) from categories c
join symptoms_categories sc on c.id = sc.category_id
join symptoms s on s.id = sc.symptom_id
join people_symptoms ps on ps.symptom_id = s.id
join vaccines v on v.id = ps.vaccine_id
where v.manufacturer = 'pfizer'
group by c.name;
*/

/* results page
table with: date of report, symptom, notes, for vaccine X, category Y, age Z, sex T

select p.reported_at, p.notes, JSONARRAYTHING(s.name) from people p
join people_symptoms ps on p.vaers_id = ps.vaers_id
join symptoms s on s.id = ps.symptom_id
join symptoms_categories sc on sc.symptom_id = s.id
join categories c on c.id = sc.category_id
join vaccines v on v.id = ps.vaccine_id
where p.sex = X and p.age between Y and Z and v.manufacturer = F and c.name = K
order by p.age;

 */

/*
how many people experienced symptoms from category X, group by vaccine manufacturer?
how many men experienced symptoms from category X, group by vaccine manufacturer?
how many women between 30 and 50 experienced symptoms from category X, group by vaccine manufacturer?

how many people experienced symptom Y from category X, group by vaccine manufacturer?
give me the notes for all the men between 60 and 70 who experienced symptom Z, from vaccine manufacturer X.
 */
