// Package sessn validades, normalizes and generates Swedish SSNs.
// It supports regular, coordination and interim numbers,
// and is written with clarity and speed in mind.
package sessn

import (
	"errors"
	"fmt"
)

type checkOptions struct {
	allowInterim      bool
	allowCoordination bool
}

// CheckOption is an option for running checks.
type CheckOption func(p *checkOptions)

// AllowInterim allows interim numbers when doing checks.
func AllowInterim(p *checkOptions) {
	p.allowInterim = true
}

// AllowCoordination allows coordination numbers when doing checks.
func AllowCoordination(p *checkOptions) {
	p.allowCoordination = true
}

// AllowAll allows both interim and coordination numbers when doing checks.
func AllowAll(p *checkOptions) {
	p.allowInterim = true
	p.allowCoordination = true
}

const (
	minYear = 1800
	maxYear = 2200
)

var (
	errBadLength      = errors.New("bad length")
	errBadFormat      = errors.New("bad format")
	errBadChecksum    = errors.New("bad checksum")
	errNoCoordination = errors.New("coordination numbers not allowed")
	errNoInterim      = errors.New("interim numbers not allowed")

	rule3 = []int{0, 2, 4, 6, 8, 1, 3, 5, 7, 9}
)

// luhn tests the validity of a number using the Luhn algorithm.
func luhn(number string) bool {
	bytes := []byte(number)
	odd := len(bytes) & 1
	var sum int
	for i, c := range bytes {
		if c < '0' || c > '9' {
			return false
		} else if i&1 == odd {
			sum += rule3[c-'0']
		} else {
			sum += int(c - '0')
		}
	}

	return sum%10 == 0
}

// Check returns errors associated with the parsed SSN, or nil if it's valid.
func (p Parsed) Check(options ...CheckOption) error {
	var opt checkOptions
	for _, setopt := range options {
		setopt(&opt)
	}

	isInterim := (p.Interim != 0)
	if p.Coordination && isInterim {
		return errBadFormat
	} else if p.Coordination && !opt.allowCoordination {
		return errNoCoordination
	} else if isInterim && !opt.allowInterim {
		return errNoInterim
	}

	if p.Year < minYear || p.Year > maxYear ||
		p.Month < 1 || p.Month > 12 ||
		p.Day < 1 || p.Day > 31 ||
		p.CheckDigit < 0 || p.CheckDigit > 9 ||
		p.BirthNumber == 0 {
		return errBadFormat
	}

	num := p.BirthNumber
	if isInterim {
		num += 100
	}

	day := p.Day
	if p.Coordination {
		day += coordinationOffset
	}

	digits := fmt.Sprintf("%04d%02d%02d%03d%d", p.Year, p.Month, day, num, p.CheckDigit)
	if tendigits := digits[2:]; !luhn(tendigits) {
		return errBadChecksum
	}

	return nil
}
