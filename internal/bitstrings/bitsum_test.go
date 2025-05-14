package bitstrings_test

import (
	"fmt"
	"testing"

	"github.com/ChainsAre2Tight/mgm-go/internal/bitstrings"
)

func TestBitSum(t *testing.T) {
	tt := []struct {
		a string
		b string
		c string
	}{
		{"1010", "1010", "0000"},
		{"110", "101", "011"},
		{"1", "01", "00"},
		{"1010", "0010", "1000"},
	}
	for _, td := range tt {
		t.Run(
			fmt.Sprintf("%s + %s = %s", td.a, td.b, td.c),
			func(t *testing.T) {
				bs_a, _ := bitstrings.FromString(td.a)
				bs_b, _ := bitstrings.FromString(td.b)
				if res, err := bitstrings.BitSum(bs_a, bs_b); err != nil || res.String() != td.c {
					t.Fatalf("\nGot:  %s, \nWant: %s, \nError: %s", res, td.c, err)
				}
			},
		)
	}
}
