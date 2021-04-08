package config

type Config struct {
	SymptomsFilePath string `env:",default=/Users/thesmurf/go/src/github.com/thehungrysmurf/vax/symptoms.csv"`
	VaccinesFilePath string `env:",default=/Users/thesmurf/go/src/github.com/thehungrysmurf/vax/vaccines.csv"`
	ReportsFilePath string `env:",default=/Users/thesmurf/go/src/github.com/thehungrysmurf/vax/reports.csv"`
	DatabaseURI string `env:",default=postgres://thesmurf:@localhost:5432/vax"`
}
