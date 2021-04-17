package config

type Config struct {
	SymptomsFilePath string `env:"SYMPTOMS_FILE_PATH,required"`
	VaccinesFilePath string `env:"VACCINES_FILE_PATH,required"`
	ReportsFilePath string `env:"REPORTS_FILE_PATH,required"`
	VaccinationTotalsFilePath string `env:"VACCINATION_TOTALS_FILE_PATH,required"`
	DatabaseURI string `env:"DB_URI,required"`
}
