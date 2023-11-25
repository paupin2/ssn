package sessn

import (
	"strconv"
	"strings"
	"time"
)

// SSN is a string wrapper used for Swedish SSNs.
type SSN string

const interimLetters = "JKLMNRSTUWX"

var (
	currentYear, _, _ = time.Now().Date()
	currentCentury    = (currentYear / 100) * 100
)

func stringHasOnlyDigits(s string) bool {
	for _, r := range s {
		if r < '0' || r > '9' {
			return false
		}
	}
	return true
}

func isInterimLetter(b byte) bool {
	return strings.ContainsRune(interimLetters, rune(b))
}

// splitDigits splits the [whitespace trimmed] SSN into its parts, or returns
// an error if the format is invalid.
func (n SSN) splitDigits() (cc, yy, mm, dd, bbb, k string, i byte, plus bool, err error) {
	parseSign := func(sign byte) bool {
		switch sign {
		case '-':
			// no effect
		case '+':
			plus = true
		default:
			err = errBadFormat
			return false
		}
		return true
	}

	// CC:century, YY:year, MM:month, DD:day, I:interim letter, BBB:birth number, K:check digit
	switch s := strings.TrimSpace(string(n)); len(s) {
	case 13:
		// CCYYMMDD-IBBK, CCYYMMDD+IBBK, CCYYMMDD-BBBK, CCYYMMDD+BBBK
		cc, yy, mm, dd, k = s[0:2], s[2:4], s[4:6], s[6:8], s[12:13]
		if !parseSign(s[8]) {
			return
		}

		if isInterimLetter(s[9]) {
			i, bbb = s[9], s[10:12] // CCYYMMDD-IBBK, CCYYMMDD+IBBK

		} else {
			bbb = s[9:12] // CCYYMMDD-BBBK, CCYYMMDD+BBBK
		}

	case 12:
		// CCYYMMDDIBBK, CCYYMMDDBBBK
		cc, yy, mm, dd, k = s[0:2], s[2:4], s[4:6], s[6:8], s[11:12]

		if isInterimLetter(s[8]) {
			i, bbb = s[8], s[9:11] // CCYYMMDDIBBK
		} else {
			bbb = s[8:11] // CCYYMMDDBBBK
		}

	case 11:
		// YYMMDD-IBBK, YYMMDD+IBBK, YYMMDD-BBBK, YYMMDD+BBBK
		yy, mm, dd, k = s[0:2], s[2:4], s[4:6], s[10:11]
		if !parseSign(s[6]) {
			return
		}

		if isInterimLetter(s[7]) {
			i, bbb = s[7], s[8:10] // YYMMDD-IBBK, YYMMDD+IBBK

		} else {
			bbb = s[7:10] // YYMMDD-BBBK, YYMMDD+BBBK
		}

	case 10:
		// YYMMDDBBBK, YYMMDDIBBK
		yy, mm, dd, k = s[0:2], s[2:4], s[4:6], s[9:10]
		if isInterimLetter(s[6]) {
			i, bbb = s[6], s[7:9] // YYMMDDIBBK

		} else {
			bbb = s[6:9] // YYMMDDBBBK
		}

	default:
		err = errBadLength
		return
	}

	if !stringHasOnlyDigits(cc + yy + mm + dd + bbb + k) {
		err = errBadFormat
	}
	return
}

func parseDigits(digits string) int {
	n, _ := strconv.Atoi(digits)
	return n
}

// Parse the SSN into its components, returning an error if the format is not
// known. This does not run consistency checks; use Parsed.Check() for that.
// The following formats are accepted (after trimming whitespace):
// "CCYYMMDD-IBBK", "CCYYMMDD+IBBK", "CCYYMMDD-BBBK", "CCYYMMDD+BBBK"
// "CCYYMMDDIBBK", "CCYYMMDDBBBK", "YYMMDD-IBBK", "YYMMDD+IBBK", "YYMMDD-BBBK",
// "YYMMDD+BBBK", "YYMMDDBBBK", "YYMMDDIBBK".
func (n SSN) Parse() (Parsed, error) {
	cc, yy, mm, dd, bbb, k, i, plus, err := n.splitDigits()
	if err != nil {
		return Parsed{}, err
	}

	return buildParsed(cc, yy, mm, dd, bbb, k, i, plus), nil
}
