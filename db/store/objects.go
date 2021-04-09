package store

import (
	"strings"
	"time"
)

type Gender string

const(
	UnknownGender Gender = "U"
	Male = "M"
	Female = "F"
)

func(g *Gender) FromString(s string) Gender {
	switch s {
	case "M":
		return Male
	case "F":
		return Female
	default:
		return UnknownGender
	}
}

type Report struct {
	VaersID int64
	Age int
	Gender Gender
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
	Covid19 Illness = "covid19"
)

type Manufacturer string

const(
	UnknownManufacturer Manufacturer = "unknown"
	Moderna = "moderna"
	Pfizer = "pfizer"
	Janssen = "janssen"
)

func(m *Manufacturer) FromString(s string) Manufacturer {
	lower := strings.ToLower(s)
	switch lower {
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

type Vaccine struct {
	Illness Illness
	Manufacturer Manufacturer
}
