package sessn

import (
	"fmt"
	"time"
)

// Gender distinguishes genders in SSNs.
type Gender int

const (
	// Female is used for SSNs belonging to females.
	Female Gender = 0
	// Male is used for SSNs belonging to males.
	Male Gender = 1

	coordinationOffset = 60
)

// Parsed represents a parsed SSN.
type Parsed struct {
	Year        int
	Month       int
	Day         int
	BirthNumber int
	CheckDigit  int
	Gender      Gender

	// Coordination numbers are unique identifiers that the Swedish Tax Agency
	// can assign to an individual who has never been listed in the Population
	// Register. No two coordination numbers are identical, and no individual
	// gets two coordination numbers.
	// These numbers have 60 added to the day field.
	Coordination bool

	// Interim numbers (T-numbers) are assigned to students who study less
	// than 12 months in Sweden and are only used in the Stockholm University.
	Interim byte
}

// Birthday returns the likely birth date associated with the SSN.
// There is the possibility, if too many birth numbers occurred at the specific
// date, that the SSN uses a different date.
func (p Parsed) Birthday() time.Time {
	return time.Date(p.Year, time.Month(p.Month), p.Day, 0, 0, 0, 0, time.UTC)
}

func (p Parsed) String() string {
	day := p.Day
	if p.Coordination {
		day += coordinationOffset
	}

	prefix := fmt.Sprintf("%04d%02d%02d", p.Year, p.Month, day)
	suffix := fmt.Sprintf("%02d%1d", p.BirthNumber, p.CheckDigit)

	if p.Interim != 0 {
		return prefix + string(p.Interim) + suffix
	}

	return prefix + suffix
}

const yearCutoff = 40

func getFullYear(yy, cc string, plus bool) int {
	year := parseDigits(yy)
	if cc != "" {
		return parseDigits(cc)*100 + year
	}

	century := currentCentury
	if century+year+yearCutoff > currentYear {
		century -= 100
	}
	if age := currentYear - (century + year); plus && age < 100 {
		century -= 100
	} else if !plus && age >= 100 {
		century += 100
	}

	return century + year
}

// buildParsed builds the parsed structure.
// It assumes all strings contain only digits.
func buildParsed(cc, yy, mm, dd, bbb, k string, i byte, plus bool) (p Parsed) {
	p.Year = getFullYear(yy, cc, plus)
	p.Month = parseDigits(mm)
	if p.Day = parseDigits(dd); p.Day > coordinationOffset {
		p.Day -= coordinationOffset
		p.Coordination = true
	}

	p.BirthNumber = parseDigits(bbb)
	p.Gender = Gender(p.BirthNumber % 2)
	p.CheckDigit = parseDigits(k)
	p.Interim = i
	return
}
