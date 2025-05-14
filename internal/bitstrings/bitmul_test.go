package bitstrings_test

import (
	"fmt"
	"testing"

	"github.com/ChainsAre2Tight/mgm-go/internal/bitstrings"
)

func TestBitstringMultiplication(t *testing.T) {
	tt := []struct {
		a string
		b string
		c string
	}{
		{"1010", "1010", "1010"},
	}
	for _, td := range tt {
		t.Run(
			fmt.Sprintf("%s x %s -> %s", td.a, td.b, td.c),
			func(t *testing.T) {
				a, _ := bitstrings.FromString(td.a)
				b, _ := bitstrings.FromString(td.b)
				if res, err := bitstrings.BitMul(a, b); err != nil || res.String() != td.c {
					t.Fatalf("\nGot:  %s, \nWant: %s, \nError: %s", res, td.c, err)
				}
			},
		)
	}
}
