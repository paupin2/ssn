package sessn

import (
	"fmt"
	"math/rand"
	"time"
)

// Valid returns whether the number is a valid Swedish SSN.
func Valid(ssn string, options ...CheckOption) bool {
	n, err := SSN(ssn).Parse()
	if err == nil {
		err = n.Check(options...)
	}
	return err == nil
}

// Normalize parses the Swedish SSN and returns its normalized (12-digit) version.
// If the number is not valid, an empty string is returned instead.
func Normalize(ssn string, options ...CheckOption) string {
	n, err := SSN(ssn).Parse()
	if err == nil {
		err = n.Check(options...)
	}
	if err == nil {
		return n.String()
	}
	return ""
}

// Parse will parse a Swedish SSN, checking format and checksum.
func Parse(ssn string, options ...CheckOption) (Parsed, error) {
	p, err := SSN(ssn).Parse()
	if err == nil {
		err = p.Check(options...)
	}
	return p, err
}

func generateForBirday(date time.Time) string {
	prefix := fmt.Sprintf("%s%03d", date.Format("20060102"), 1+rand.Int()%999)

	// brute-force checksum
	for k := '0'; k <= '9'; k++ {
		if n, err := Parse(prefix + string(k)); err == nil {
			return n.String()
		}
	}

	return "" // this will not happen
}

// Generate returns a random regular Swedish SSN.
func Generate() string {
	start := time.Date(minYear, 1, 1, 0, 0, 0, 0, time.UTC)
	days := rand.Int() % ((currentYear - minYear - 1) * 365)
	date := start.Add(24 * time.Hour * time.Duration(days))
	return generateForBirday(date)
}
