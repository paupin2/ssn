package sessn

import (
	"testing"
)

func TestGenerate(t *testing.T) {
	total, failed := 1000, 0
	for i := 0; i < total; i++ {
		if Generate() == "" {
			t.Fatalf("failed to generate %d times out of %d", failed, total)
		}
	}
}
