package config

type Config struct {
	SymptomsFilePath string `env:",default=/Users/thesmurf/go/src/github.com/thehungrysmurf/vax/test_data/symptoms.csv"`
	VaccinesFilePath string `env:",default=/Users/thesmurf/go/src/github.com/thehungrysmurf/vax/test_data/vaccines.csv"`
	ReportsFilePath string `env:",default=/Users/thesmurf/go/src/github.com/thehungrysmurf/vax/test_data/reports.csv"`
	DatabaseURI string `env:",default=postgres://thesmurf:@localhost:5432/vax"`
}
