package sessn

import (
	"compress/bzip2"
	"embed"
	"encoding/json"
	"fmt"
	"testing"
)

type testItem struct {
	Integer       int    `json:"integer"`
	Long          string `json:"long_format"`
	Short         string `json:"short_format"`
	Separated     string `json:"separated_format"`
	SeparatedLong string `json:"separated_long"`
	Valid         bool   `json:"valid"`
	Type          string `json:"type"`
	IsMale        bool   `json:"isMale"`
	IsFemale      bool   `json:"isFemale"`
}

var (
	//go:embed testdata/*.json.bz2
	testFiles embed.FS
)

func readItems(t *testing.T, fn string) (items []testItem) {
	fn = "testdata/" + fn + ".json.bz2"
	t.Helper()
	f, err := testFiles.Open(fn)
	if err != nil {
		t.Fatalf("opening %s: %v", fn, err)
	}
	defer f.Close()

	decoder := json.NewDecoder(bzip2.NewReader(f))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&items); err != nil {
		t.Fatalf("parsing %s: %v", fn, err)
	}
	return items
}

func readAllItems(t *testing.T) []testItem {
	t.Helper()
	return append(
		readItems(t, "list"),
		readItems(t, "interim")...,
	)
}

func TestSSN_Parse(t *testing.T) {
	check := func(in string, expected Parsed, expectedErr string) {
		t.Helper()
		actual, err := SSN(in).Parse()
		var actualErr string
		if err != nil {
			actualErr = err.Error()
		}

		if actualErr != expectedErr {
			t.Errorf("%q: expected error %q but got %q instead", in, expectedErr, actualErr)
		}
		if actual != expected {
			t.Errorf("%q: expected %+v but got %+v instead", in, expected, actual)
		}
	}

	const noError = ""
	check("19791207-X239", Parsed{Year: 1979, Month: 12, Day: 7, Interim: 'X', CheckDigit: 9, BirthNumber: 23, Gender: Male}, noError)
	check("19791207+X239", Parsed{Year: 1979, Month: 12, Day: 7, Interim: 'X', CheckDigit: 9, BirthNumber: 23, Gender: Male}, noError)
	check("19791207-1239", Parsed{Year: 1979, Month: 12, Day: 7, CheckDigit: 9, BirthNumber: 123, Gender: Male}, noError)
	check("19791207+1239", Parsed{Year: 1979, Month: 12, Day: 7, CheckDigit: 9, BirthNumber: 123, Gender: Male}, noError)
	check("19791207X239", Parsed{Year: 1979, Month: 12, Day: 7, Interim: 'X', CheckDigit: 9, BirthNumber: 23, Gender: Male}, noError)
	check("197912071239", Parsed{Year: 1979, Month: 12, Day: 7, CheckDigit: 9, BirthNumber: 123, Gender: Male}, noError)
	check("791207-X239", Parsed{Year: 1979, Month: 12, Day: 7, Interim: 'X', CheckDigit: 9, BirthNumber: 23, Gender: Male}, noError)
	check("791207+X239", Parsed{Year: 1879, Month: 12, Day: 7, Interim: 'X', CheckDigit: 9, BirthNumber: 23, Gender: Male}, noError)
	check("791207-1239", Parsed{Year: 1979, Month: 12, Day: 7, CheckDigit: 9, BirthNumber: 123, Gender: Male}, noError)
	check("791207+1239", Parsed{Year: 1879, Month: 12, Day: 7, CheckDigit: 9, BirthNumber: 123, Gender: Male}, noError)
	check("7912071239", Parsed{Year: 1979, Month: 12, Day: 7, CheckDigit: 9, BirthNumber: 123, Gender: Male}, noError)
	check("791207X239", Parsed{Year: 1979, Month: 12, Day: 7, Interim: 'X', CheckDigit: 9, BirthNumber: 23, Gender: Male}, noError)

	check(" 19791207-X239", Parsed{Year: 1979, Month: 12, Day: 7, Interim: 'X', CheckDigit: 9, BirthNumber: 23, Gender: Male}, noError)
	check("19791207-X239 ", Parsed{Year: 1979, Month: 12, Day: 7, Interim: 'X', CheckDigit: 9, BirthNumber: 23, Gender: Male}, noError)
	check("   19791207-X239   ", Parsed{Year: 1979, Month: 12, Day: 7, Interim: 'X', CheckDigit: 9, BirthNumber: 23, Gender: Male}, noError)

	check("19791267-X239", Parsed{Year: 1979, Month: 12, Day: 7, Interim: 'X', CheckDigit: 9, BirthNumber: 23, Coordination: true, Gender: Male}, noError)
	check("19791267+X239", Parsed{Year: 1979, Month: 12, Day: 7, Interim: 'X', CheckDigit: 9, BirthNumber: 23, Coordination: true, Gender: Male}, noError)
	check("19791267-1239", Parsed{Year: 1979, Month: 12, Day: 7, CheckDigit: 9, BirthNumber: 123, Coordination: true, Gender: Male}, noError)
	check("19791267+1239", Parsed{Year: 1979, Month: 12, Day: 7, CheckDigit: 9, BirthNumber: 123, Coordination: true, Gender: Male}, noError)
	check("19791267X239", Parsed{Year: 1979, Month: 12, Day: 7, Interim: 'X', CheckDigit: 9, BirthNumber: 23, Coordination: true, Gender: Male}, noError)
	check("197912671239", Parsed{Year: 1979, Month: 12, Day: 7, CheckDigit: 9, BirthNumber: 123, Coordination: true, Gender: Male}, noError)
	check("791267-X239", Parsed{Year: 1979, Month: 12, Day: 7, Interim: 'X', CheckDigit: 9, BirthNumber: 23, Coordination: true, Gender: Male}, noError)
	check("791267+X239", Parsed{Year: 1879, Month: 12, Day: 7, Interim: 'X', CheckDigit: 9, BirthNumber: 23, Coordination: true, Gender: Male}, noError)
	check("791267-1239", Parsed{Year: 1979, Month: 12, Day: 7, CheckDigit: 9, BirthNumber: 123, Coordination: true, Gender: Male}, noError)
	check("791267+1239", Parsed{Year: 1879, Month: 12, Day: 7, CheckDigit: 9, BirthNumber: 123, Coordination: true, Gender: Male}, noError)
	check("7912671239", Parsed{Year: 1979, Month: 12, Day: 7, CheckDigit: 9, BirthNumber: 123, Coordination: true, Gender: Male}, noError)
	check("791267X239", Parsed{Year: 1979, Month: 12, Day: 7, Interim: 'X', CheckDigit: 9, BirthNumber: 23, Coordination: true, Gender: Male}, noError)

	check("", Parsed{}, "bad length")
	check("1", Parsed{}, "bad length")
	check("12", Parsed{}, "bad length")
	check("123", Parsed{}, "bad length")
	check("1234", Parsed{}, "bad length")
	check("12345", Parsed{}, "bad length")
	check("123456", Parsed{}, "bad length")
	check("1234567", Parsed{}, "bad length")
	check("12345678", Parsed{}, "bad length")
	check("123456789", Parsed{}, "bad length")
	check("A123456789", Parsed{}, "bad format")
	check("AB123456789", Parsed{}, "bad format")
	check("ABC123456789", Parsed{}, "bad format")
	check("ABCD123456789", Parsed{}, "bad format")
	check("ABCDE123456789", Parsed{}, "bad length")

	check("00000000/0000", Parsed{}, "bad format")
	check("000000/000", Parsed{}, "bad format")
}

func equal[T comparable](t *testing.T, actual, expected T, msg string, a ...any) {
	t.Helper()
	if actual != expected {
		t.Errorf("%s: expected %v but got %v", fmt.Sprintf(msg, a...), expected, actual)
	}
}

func TestSSN_Parse2(t *testing.T) {
	for _, expected := range readAllItems(t) {
		// parse all variations
		for _, s := range []string{expected.Long, expected.Short, expected.Separated, expected.SeparatedLong} {
			if _, err := SSN(s).Parse(); err != nil && expected.Valid {
				t.Errorf("expected %q to be valid but got %v", s, err)
			}
		}

		in := expected.Long
		p, err := SSN(in).Parse()
		if err != nil && expected.Valid {
			t.Errorf("expected %q to be valid but got %v", in, err)
		}
		if !expected.Valid {
			return
		}

		equal(t, p.String(), expected.Long, "%q: long", in)
		if actual := (p.Gender == Male); expected.IsMale != actual {
			t.Errorf("on %q: expected male? %t but got %t", in, expected.IsMale, actual)
		}
		if actual := (p.Gender == Female); expected.IsFemale != actual {
			t.Errorf("on %q: expected female? %t but got %t", in, expected.IsFemale, actual)
		}
		if expected.Type == "con" && !p.Coordination {
			t.Errorf("on %q: expected coordination", in)
		}
		if expected.Type == "interim" && p.Interim == 0 {
			t.Errorf("on %q: expected interim", in)
		}
		if expected.Type == "ssn" && (p.Coordination || p.Interim != 0) {
			t.Errorf("on %q: expected regular ssn", in)
		}
	}
}
