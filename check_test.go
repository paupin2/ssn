package sessn

import (
	"bufio"
	"compress/bzip2"
	"path"
	"testing"
)

func TestParsed_Check(t *testing.T) {
	for _, item := range readAllItems(t) {
		for _, s := range []string{item.Long, item.Short, item.Separated, item.SeparatedLong} {
			p, err := SSN(s).Parse()
			if err == nil {
				err = p.Check(AllowAll)
			}
			if valid := (err == nil); item.Valid != valid {
				t.Errorf("on %q expected valid? %t %t: %v", s, item.Valid, valid, err)
			}
		}
	}
}

func forEachTestNumber(t *testing.T, check func(s string) error) {
	t.Helper()
	entries, err := testFiles.ReadDir("testdata")
	if err != nil {
		t.Fatalf("opening dir: %v", err)
	}
	for _, e := range entries {
		fn := e.Name()
		if !e.IsDir() && path.Ext(fn) == ".csv.bz2" {
			f, err := testFiles.Open(fn)
			if err != nil {
				t.Fatalf("opening %s: %v", fn, err)
			}
			defer f.Close()

			scanner := bufio.NewScanner(bzip2.NewReader(f))
			for scanner.Scan() {
				if number := scanner.Text(); number != "" {
					if err := check(number); err != nil {
						t.Errorf("parsing %q: %v", number, err)
					}
				}
			}

		}
	}
}

func TestParsed_CheckNumbers(t *testing.T) {
	forEachTestNumber(t, func(s string) error {
		p, err := SSN(s).Parse()
		if err == nil {
			err = p.Check(AllowAll)
		}
		if err != nil {
			t.Errorf("parsing %q: %v", s, err)
		}
		return err
	})
}
