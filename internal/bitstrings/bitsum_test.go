package bitstrings_test

import (
	"fmt"
	"testing"

	"github.com/ChainsAre2Tight/mgm-go/internal/bitstrings"
)

func TestBitSumFromString(t *testing.T) {
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
				if res := bitstrings.BitSum128(bs_a, bs_b); res.String() != td.c {
					t.Fatalf("\nGot:  %s, \nWant: %s.", res, td.c)
				}
			},
		)
	}
}

func TestBitSum(t *testing.T) {
	tt := []struct {
		a, b, c *bitstrings.BitString128
	}{
		{
			a: bitstrings.FromGOSTString("B8 57 48 C5 12 F3 19 90 AA 56 7E F1 53 35 DB 74"),
			b: bitstrings.FromGOSTString("11 22 33 44 55 66 77 00 FF EE DD CC BB AA 99 88"),
			c: bitstrings.FromGOSTString("A9 75 7B 81 47 95 6E 90 55 B8 A3 3D E8 9F 42 FC"),
		},
	}
	for _, td := range tt {
		t.Run(
			fmt.Sprintf("%x | %x + %x | %x = %x | %x", td.a.Upper(), td.a.Lower(), td.b.Upper(), td.b.Lower(), td.c.Upper(), td.c.Lower()),
			func(t *testing.T) {
				if res := bitstrings.BitSum128(td.a, td.b); res.Upper() != td.c.Upper() || res.Lower() != td.c.Lower() {
					t.Fatalf("Got:  %x , %x, \nWant: %x , %x.", res.Upper(), res.Lower(), td.c.Upper(), td.c.Lower())
				}
			},
		)
	}
}
