package store

import (
	"strings"
	"time"
)

type Sex string

const(
	UnknownSex Sex = "U"
	Male = "M"
	Female = "F"
)

func(s *Sex) FromString(str string) Sex {
	switch str {
	case "M":
		return Male
	case "F":
		return Female
	default:
		return UnknownSex
	}
}

func(s *Sex) String() string {
	switch *s {
	case Male:
		return "Male"
	case Female:
		return "Female"
	default:
		return "Unknown"
	}
}

type Report struct {
	VaersID int64
	Age int
	Sex Sex
	Notes string
	ReportedAt time.Time
}

type Symptom struct {
	ID int64
	Name string
	Alias string
	CategoryIDs []int
}

type Illness string

const(
	UnknownIllness Illness = "unknown"
	Covid19 = "covid19"
)

type Manufacturer string

const(
	UnknownManufacturer Manufacturer = "unknown"
	Moderna = "moderna"
	Pfizer = "pfizer"
	Janssen = "janssen"
)

func(m *Manufacturer) FromString(s string) Manufacturer {
	switch strings.ToLower(s) {
		case "moderna":
			return Moderna
		case `pfizer\biontech`:
			return Pfizer
		case "janssen":
			return Janssen
		default:
			return UnknownManufacturer
	}
}

func(m *Manufacturer) String() string {
	switch *m {
	case Moderna:
		return "Moderna"
	case Pfizer:
		return "Pfizer"
	case Janssen:
		return "Johnson & Johnson"
	default:
		return "Unknown"
	}
}

type Vaccine struct {
	Illness Illness
	Manufacturer Manufacturer
}

func(i *Illness) FromString(s string) Illness {
	switch strings.ToLower(s) {
	case "covid19":
		return Covid19
	default:
		return UnknownIllness
	}
}
