package config

type Config struct {
	SymptomsFilePath string `env:",default=/Users/thesmurf/Documents/2021VAERSData/2021VAERSSYMPTOMS.csv"`
	VaccinesFilePath string `env:",default=/Users/thesmurf/Documents/2021VAERSData/2021VAERSVAX.csv"`
	ReportsFilePath string `env:",default=/Users/thesmurf/Documents/2021VAERSData/2021VAERSDATA.csv"`
	DatabaseURI string `env:",default=postgres://thesmurf:@localhost:5432/vax"`
}
